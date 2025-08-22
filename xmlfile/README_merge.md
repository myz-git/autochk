# XML文件合并工具

## 功能说明

`merge.go` 是一个智能的XML文件合并工具，专门用于处理Oracle RAC环境的健康检查报告文件。

### 主要功能：

1. **自动扫描**：扫描 `input_xml/` 目录下的所有XML文件
2. **智能识别**：自动识别单实例文件和RAC文件
3. **RAC合并**：按照规则合并RAC文件
4. **单实例处理**：直接复制单实例文件到输出目录

## 使用方法

### 1. 编译程序
```bash
go build -o merge.exe merge.go
```

### 2. 准备文件
将XML文件放入 `input_xml/` 目录，文件命名格式：
```
YYYYMMDD_HOSTNAME_DBNAME_INSTNAME.xml
```

### 3. 运行程序
```bash
./merge.exe
```

### 4. 查看结果
合并后的文件将保存在 `output_xml/` 目录中

## 文件命名规则

### 输入文件格式
- 格式：`YYYYMMDD_HOSTNAME_DBNAME_INSTNAME.xml`
- 示例：
  - `20250606_myzrac11_racdb_racdb1.xml`
  - `20250606_myzrac12_racdb_racdb2.xml`
  - `20240401_m0ora01_M001R01DG_M001R011.xml`

### 输出文件格式
- **单实例文件**：保持原文件名不变
- **RAC合并文件**：`YYYYMMDD_HOSTNAME1.HOSTNAME2_DBNAME_RAC.xml`
- 示例：
  - `20250606_myzrac11.myzrac12_racdb_RAC.xml`
  - `20240401_m0ora01.lm0ora02.lm0ora03_M001R01DG_RAC.xml`

## 合并规则

### RAC文件识别
程序通过以下规则识别RAC文件：
1. 相同日期（YYYYMMDD）
2. 相同数据库名（DBNAME）
3. 不同主机名（HOSTNAME）
4. 实例名最后两位数字不同

### 合并逻辑
1. **以实例号最小的文件为主文件**
2. **TAG0处理**：将其他文件的TAG0/NODE1内容复制为主文件的TAG0/NODE2、NODE3等
3. **TAG1处理**：保持不变（数据库级别信息相同）
4. **TAG2处理**：将其他文件的TAG2/NODE1内容复制为主文件的TAG2/NODE2、NODE3等

### 示例
**输入文件：**
- `20250606_myzrac11_racdb_racdb1.xml` (实例号1)
- `20250606_myzrac12_racdb_racdb2.xml` (实例号2)

**合并后结构：**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<VER></VER>
<EACHK>
    <TAG0>
        <NODE1>
            <!-- myzrac11的主机信息 -->
        </NODE1>
        <NODE2>
            <!-- myzrac12的主机信息 -->
        </NODE2>
    </TAG0>
    <TAG1>
        <NODE1>
            <!-- 数据库信息（保持不变） -->
        </NODE1>
    </TAG1>
    <TAG2>
        <NODE1>
            <!-- racdb1实例信息 -->
        </NODE1>
        <NODE2>
            <!-- racdb2实例信息 -->
        </NODE2>
    </TAG2>
</EACHK>
```

## 目录结构

```
项目根目录/
├── merge.go
├── merge.exe
├── input_xml/          # 输入XML文件目录
│   ├── 20250606_myzrac11_racdb_racdb1.xml
│   ├── 20250606_myzrac12_racdb_racdb2.xml
│   └── ...
└── output_xml/         # 输出文件目录
    ├── 20250606_myzrac11.myzrac12_racdb_RAC.xml
    └── ...
```

## 错误处理

- 自动跳过无法读取或解析的文件
- 记录详细的处理日志
- 继续处理其他文件
- 单实例文件直接复制，不进行合并

## 注意事项

1. 确保输入目录 `input_xml/` 存在
2. 程序会自动创建输出目录 `output_xml/`
3. 文件名必须严格按照格式：`YYYYMMDD_HOSTNAME_DBNAME_INSTNAME.xml`
4. 实例名必须包含数字后缀用于排序
5. 使用 `github.com/beevik/etree` 库处理XML，确保已安装依赖

## 依赖安装

```bash
go mod init autochk
go get github.com/beevik/etree
``` 