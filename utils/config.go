package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var configFile []byte
var ColInx int

// 解析rule.yaml文件
// //*** Lv1***////
type RuleInfo struct {
	Osrule Osrule `yaml:"osrule"`
	Dbrule Dbrule `yaml:"dbrule"`
}

// //*** Lv2***////
type Osrule struct {
	Osparameter Osparameter `yaml:"osparameter"`
	Ulimit      Ulimit      `yaml:"ulimit"`
	Filesystem  Filesystem  `yaml:"filesystem"`
	Inodeusage  Inodeusage  `yaml:"inodeusage"`
	Cpustat     Cpustat     `yaml:"cpustat"`
	Memstat     Memstat     `yaml:"memstat"`
	Iostat      Iostat      `yaml:"iostat"`
	Thpstat     Thpstat     `yaml:"thpstat"`
	Numa        Numa        `yaml:"numa"`
	Ntp         Ntp         `yaml:"ntp"`
}

// //*** Lv2***////
type Dbrule struct {
	// 数据库实例分析
	Dbstatus          Dbstatus          `yaml:"dbstatus"`
	Dbtbsusage        Dbtbsusage        `yaml:"dbtbsusage"`
	Dbdatafile        Dbdatafile        `yaml:"dbdatafile"`
	Dbcontrolfile     Dbcontrolfile     `yaml:"dbcontrolfile"`
	Dbusersize        Dbusersize        `yaml:"dbusersize"`
	Dbredocheck       Dbredocheck       `yaml:"dbredocheck"`
	Dbredoswitch      Dbredoswitch      `yaml:"dbredoswitch"`
	Dbparameter       Dbparameter       `yaml:"dbparameter"`
	Db_parameter_file Db_parameter_file `yaml:"db_parameter_file"`
	Db_shp_size       Db_shp_size       `yaml:"db_shp_size"`
	Db_shp_pct        Db_shp_pct        `yaml:"db_shp_pct"`

	// 数据库对象分析
	Dbtableparallel Dbtableparallel `yaml:"dbtableparallel"`
	Dbindexparallel Dbindexparallel `yaml:"dbindexparallel"`
	Dbinvalidindex  Dbinvalidindex  `yaml:"dbinvalidindex"`
	Dbsequence      Dbsequence      `yaml:"dbsequence"`
	Db_seq_usage    Db_seq_usage    `yaml:"db_seq_usage"`

	// 数据库性能分析
	Db_4031check      Db_4031check      `yaml:"db_4031check"`
	Dbresource        Dbresource        `yaml:"dbresource"`
	Loadprofile       Loadprofile       `yaml:"loadprofile"`
	Instefficiency    Instefficiency    `yaml:"instefficiency"`
	Topevent          Topevent          `yaml:"topevent"`
	Topsqlbyelapstime Topsqlbyelapstime `yaml:"topsqlbyelapstime"`

	// 数据库安全检查
	Db_expir_user            Db_expir_user            `yaml:"db_expir_user"`
	Dbproductuserfailedlogin Dbproductuserfailedlogin `yaml:"dbproductuserfailedlogin"`
	Dbdbapriv                Dbdbapriv                `yaml:"dbdbapriv"`
	Dbsysdba                 Dbsysdba                 `yaml:"dbsysdba"`
	Dbauditsegment           Dbauditsegment           `yaml:"dbauditsegment"`
	Dbauditcont              Dbauditcont              `yaml:"dbauditcont"`
	Db_Nosys_In_System       Db_Nosys_In_System       `yaml:"db_nosys_in_system"`
	Dbvirscheck              Dbvirscheck              `yaml:"dbvirscheck"`
	Dbrmancheck              Dbrmancheck              `yaml:"dbrmancheck"`
	Dbscnhealthcheck         Dbscnhealthcheck         `yaml:"dbscnhealthcheck"`

	// 数据库监控、DataGuard、备份及杂项分析
	Dberrlog              Dberrlog              `yaml:"dberrlog"`
	Dbdglagcheck          Dbdglagcheck          `yaml:"dbdglagcheck"`
	Dbdgerrcheck          Dbdgerrcheck          `yaml:"dbdgerrcheck"`
	Dbcrscheck            Dbcrscheck            `yaml:"dbcrscheck"`
	Dbasmusage            Dbasmusage            `yaml:"dbasmusage"`
	Dblsnrinfo            Dblsnrinfo            `yaml:"dblsnrinfo"`
	Dbrecoverydest        Dbrecoverydest        `yaml:"dbrecoverydest"`
	Dbflashrecoveryuseage Dbflashrecoveryuseage `yaml:"dbflashrecoveryuseage"`
	Dbpsu                 Dbpsu                 `yaml:"dbpsu"`
	Dbpatch               Dbpatch               `yaml:"dbpatch"`
}

// ///*** Lv3 OS Start***/////
type Osparameter struct {
	Nm                        string   `yaml:"nm"`
	Title                     string   `yaml:"title"`
	Desc                      string   `yaml:"desc"`
	L_nproc_ne                int      `yaml:"l_nproc_ne"`
	L_nofile_ne               int      `yaml:"l_nofile_ne"`
	L_randomize_va_space      int      `yaml:"l_randomize_va_space"`
	L_panic_on_oops           int      `yaml:"l_panic_on_oops"`
	L_min_free_kbytes         int      `yaml:"l_min_free_kbytes"`
	S_disable_ism_large_pages []string `yaml:"s_disable_ism_large_pages,flow"` //返回字符串数组, flow为固定词
}

type Ulimit struct {
	Nm                   string `yaml:"nm"`
	Title                string `yaml:"title"`
	Desc                 string `yaml:"desc"`
	Open_files_ne        int    `yaml:"open_files_ne"`
	Max_user_rocesses_ne int    `yaml:"max_user_rocesses_ne"`
}

type Filesystem struct {
	Nm       string `yaml:"nm"`
	Title    string `yaml:"title"`
	Desc     string `yaml:"desc"`
	Disk_ge1 string `yaml:"disk_ge1"`
	Disk_ge2 string `yaml:"disk_ge2"`
}

type Inodeusage struct {
	Nm        string `yaml:"nm"`
	Title     string `yaml:"title"`
	Desc      string `yaml:"desc"`
	Inode_ge1 string `yaml:"inode_ge1"`
	Inode_ge2 string `yaml:"inode_ge2"`
}

type Cpustat struct {
	Nm       string `yaml:"nm"`
	Title    string `yaml:"title"`
	Desc     string `yaml:"desc"`
	Idle_le1 int    `yaml:"idle_le1"`
	Idle_le2 int    `yaml:"idle_le2"`
	Swap_ge1 int    `yaml:"swap_ge1"`
	Swap_ge2 int    `yaml:"swap_ge2"`
}
type Memstat struct {
	Nm            string `yaml:"nm"`
	Title         string `yaml:"title"`
	Desc          string `yaml:"desc"`
	Available_le1 int    `yaml:"available_le1"`
	Available_le2 int    `yaml:"available_le2"`
}

type Iostat struct {
	Nm           string  `yaml:"nm"`
	Title        string  `yaml:"title"`
	Desc         string  `yaml:"desc"`
	Diskutil_ge1 float64 `yaml:"diskutil_ge1"`
	Diskutil_ge2 float64 `yaml:"diskutil_ge2"`
}

type Thpstat struct {
	Nm         string `yaml:"nm"`
	Title      string `yaml:"title"`
	Desc       string `yaml:"desc"`
	Anpages_gt int    `yaml:"anpages_gt"`
}

type Numa struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Flg1  string `yaml:"flg1"`
	Flg2  string `yaml:"flg2"`
}

type Ntp struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

/////*** Lv3 os End***/////

// /////解析rule.yaml中的数据库部份规则
// ///*** Lv3 db Start***/////

type Dbstatus struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbtbsusage struct {
	Nm           string  `yaml:"nm"`
	Title        string  `yaml:"title"`
	Desc         string  `yaml:"desc"`
	Tbsutil_ge1  float64 `yaml:"tbsutil_ge1"`
	Tbsutil_ge2  float64 `yaml:"tbsutil_ge2"`
	Freesize_le1 float64 `yaml:"freesize_le1"`
	Freesize_le2 float64 `yaml:"freesize_le2"`
}

type Dbdatafile struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Status string `yaml:"status"`
}

type Dbcontrolfile struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Cnt_le1 int    `yaml:"cnt_le1"`
}
type Dbredocheck struct {
	Nm           string  `yaml:"nm"`
	Title        string  `yaml:"title"`
	Desc         string  `yaml:"desc"`
	Rdf_size_lt1 float64 `yaml:"rdf_size_lt1"`
	// Rdf_status   string `yaml:"rdf_status"`
	Rdf_status_list []string `yaml:"rdf_status_list,flow"` //返回字符串数组, flow为固定词
}

type Dbredoswitch struct {
	Nm         string `yaml:"nm"`
	Title      string `yaml:"title"`
	Desc       string `yaml:"desc"`
	Sw_cnt_ge1 int    `yaml:"sw_cnt_ge1"`
}

type Dbresource struct {
	Nm          string `yaml:"nm"`
	Title       string `yaml:"title"`
	Desc        string `yaml:"desc"`
	Res_use_ge1 int    `yaml:"res_use_ge1"`
}

type Loadprofile struct {
	Nm           string  `yaml:"nm"`
	Title        string  `yaml:"title"`
	Desc         string  `yaml:"desc"`
	Redosize_ge1 float64 `yaml:"redosize_ge1"`
	Logons_ge1   float64 `yaml:"logons_ge1"`
}

type Instefficiency struct {
	Nm          string  `yaml:"nm"`
	Title       string  `yaml:"title"`
	Desc        string  `yaml:"desc"`
	Buffer_hit  float64 `yaml:"buffer_hit"`
	Library_hit float64 `yaml:"library_hit"`
	Soft_parse  float64 `yaml:"soft_parse"`
}

type Topevent struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Topsqlbyelapstime struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dblsnrinfo struct {
	Nm       string `yaml:"nm"`
	Title    string `yaml:"title"`
	Desc     string `yaml:"desc"`
	Log_size int    `yaml:"log_size"`
}

type Dbtableparallel struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbindexparallel struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbinvalidindex struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbsequence struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Db_seq_usage struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbrecoverydest struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbflashrecoveryuseage struct {
	Nm      string  `yaml:"nm"`
	Title   string  `yaml:"title"`
	Desc    string  `yaml:"desc"`
	Useage1 float64 `yaml:"useage1"`
	Useage2 float64 `yaml:"useage2"`
}

type Dberrlog struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB string `yaml:"resultB"`
}

type Db_expir_user struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbproductuserfailedlogin struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbdbapriv struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultG string `yaml:"resultG"`
}

type Dbsysdba struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB string `yaml:"resultB"`
}

type Dbdglagcheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB int    `yaml:"resultB"`
}

type Dbdgerrcheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultG string `yaml:"resultG"`
}

type Dbauditsegment struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultG string `yaml:"resultG"`
}

type Dbauditcont struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultG int    `yaml:"resultG"`
}

type Db_Nosys_In_System struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB string `yaml:"resultB"`
}

type Dbvirscheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultR string `yaml:"resultR"`
}

type Dbscnhealthcheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Resulta string `yaml:"resulta"`
	Resultb string `yaml:"resultb"`
	Resultc string `yaml:"resultc"`
}

type Dbusersize struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbparameter struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Db_parameter_file struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Db_shp_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Db_shp_pct struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Db_4031check struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbcrscheck struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbasmusage struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbrmancheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB string `yaml:"resultB"`
	ResultR string `yaml:"resultR"`
}

type Dbpsu struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbpatch struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

/////*** Lv3 db End***/////

func GetRule() (c *RuleInfo, err error) {
	err = yaml.Unmarshal(configFile, &c)
	return c, err
}

func init() {
	var err error
	// configFile, err = ioutil.ReadFile("./rule.yaml")
	configFile, err = os.ReadFile("./rule.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
}
