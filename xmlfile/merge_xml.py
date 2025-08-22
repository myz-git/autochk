import glob
import re
import os
from collections import defaultdict
import xml.etree.ElementTree as ET
from typing import List, Tuple

def parse_filename(filename: str) -> Tuple[str, str, str, str]:
    """解析XML文件名，提取日期、主机名、数据库名、实例名"""
    basename = os.path.basename(filename)
    # 匹配命名规则：yyyymmdd_hostname_dbname_instname.xml
    pattern = r'(\d{8})_([^_]+)_([^_]+)_([^_]+)\.xml'
    match = re.match(pattern, basename)
    if not match:
        raise ValueError(f"文件名 {basename} 不符合命名规则")
    yyyymmdd, hostname, dbname, instname = match.groups()
    return yyyymmdd, hostname, dbname, instname

def get_instance_number(instname: str) -> int:
    """提取实例名最右侧两位的数字序号"""
    # 提取最后两位数字
    match = re.search(r'\d{1}$', instname)
    if not match:
        raise ValueError(f"实例名 {instname} 没有以两位数字结尾")
    return int(match.group())

def group_files(xml_files: List[str]) -> List[List[str]]:
    """按日期和数据库名分组文件，返回需要合并的文件组"""
    groups = defaultdict(list)
    for xml_file in xml_files:
        yyyymmdd, hostname, dbname, instname = parse_filename(xml_file)
        key = (yyyymmdd, dbname)
        groups[key].append((xml_file, hostname, instname))
    
    merge_groups = []
    for key, files in groups.items():
        if len(files) > 1:  # 只有多于一个文件才需要合并
            # 检查实例名序号是否不同
            instance_numbers = {get_instance_number(f[2]) for f in files}
            if len(instance_numbers) == len(files):  # 实例名序号各不相同
                merge_groups.append([f[0] for f in sorted(files, key=lambda x: get_instance_number(x[2]))])
    return merge_groups

def merge_xml_files(file_group: List[str], output_dir: str) -> str:
    """合并一组XML文件"""
    # 解析主文件（实例名序号最小的文件）
    main_file = file_group[0]
    yyyymmdd, main_hostname, dbname, _ = parse_filename(main_file)
    
    # 构造输出文件名
    hostnames = [parse_filename(f)[1] for f in file_group]
    output_filename = f"{yyyymmdd}_{'.'.join(hostnames)}_{dbname}_RAC.xml"
    output_path = os.path.join(output_dir, output_filename)
    
    # 解析主XML
    main_tree = ET.parse(main_file)
    main_root = main_tree.getroot()
    main_tag0 = main_root.find('.//TAG0')
    main_tag2 = main_root.find('.//TAG2')
    
    # 合并其他文件的TAG0和TAG2
    for i, xml_file in enumerate(file_group[1:], start=2):
        tree = ET.parse(xml_file)
        root = tree.getroot()
        
        # 合并TAG0
        node = root.find('.//TAG0/NODE1')
        if node is not None:
            node.tag = f'NODE{i}'
            main_tag0.append(node)
        
        # 合并TAG2
        node = root.find('.//TAG2/NODE1')
        if node is not None:
            node.tag = f'NODE{i}'
            main_tag2.append(node)
    
    # 保存合并后的XML
    main_tree.write(output_path, encoding='UTF-8', xml_declaration=True)
    return output_path

def main(input_dir: str, output_dir: str):
    """主函数，处理所有XML文件"""
    # 确保输出目录存在
    os.makedirs(output_dir, exist_ok=True)
    
    # 获取所有XML文件
    xml_files = glob.glob(os.path.join(input_dir, '*.xml'))
    if not xml_files:
        print("未找到XML文件")
        return
    
    # 分组需要合并的文件
    merge_groups = group_files(xml_files)
    
    # 合并每组文件
    for group in merge_groups:
        output_file = merge_xml_files(group, output_dir)
        print(f"已生成合并文件: {output_file}")

if __name__ == '__main__':
    input_dir = 'input_xml'  # 输入XML文件目录
    output_dir = 'output_xml'  # 输出合并文件目录
    main(input_dir, output_dir)