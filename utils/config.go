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
	Cpustat     Cpustat     `yaml:"cpustat"`
	Memstat     Memstat     `yaml:"memstat"`
	Iostat      Iostat      `yaml:"iostat"`
	Thpstat     Thpstat     `yaml:"thpstat"`
	Numa        Numa        `yaml:"numa"`
}

// //*** Lv2***////
type Dbrule struct {
	DbTbsusage               DbTbsusage               `yaml:"dbtbsusage"`
	Dbdatafile               Dbdatafile               `yaml:"dbdatafile"`
	Dbcontrolfile            Dbcontrolfile            `yaml:"dbcontrolfile"`
	Dbredocheck              Dbredocheck              `yaml:"dbredocheck"`
	Dbredoswitch             Dbredoswitch             `yaml:"dbredoswitch"`
	Dbresource               Dbresource               `yaml:"dbresource"`
	Loadprofile              Loadprofile              `yaml:"loadprofile"`
	Instefficiency           Instefficiency           `yaml:"instefficiency"`
	Dblsnrinfo               Dblsnrinfo               `yaml:"dblsnrinfo"`
	Dbtableparallel          Dbtableparallel          `yaml:"dbtableparallel"`
	Dbindexparallel          Dbindexparallel          `yaml:"dbindexparallel"`
	Dbinvalidindex           Dbinvalidindex           `yaml:"dbinvalidindex"`
	Dbsequence               Dbsequence               `yaml:"dbsequence"`
	Dbrecoverydest           Dbrecoverydest           `yaml:"dbrecoverydest"`
	Dbflashrecoveryuseage    Dbflashrecoveryuseage    `yaml:"dbflashrecoveryuseage"`
	Dberrlog                 Dberrlog                 `yaml:"dberrlog"`
	Dbproductuserfailedlogin Dbproductuserfailedlogin `yaml:"dbproductuserfailedlogin"`
	Dbdglagcheck             Dbdglagcheck             `yaml:"dbdglagcheck"`
	Dbdgerrcheck             Dbdgerrcheck             `yaml:"dbdgerrcheck"`
	Dbrmancheck              Dbrmancheck              `yaml:"dbrmancheck"`
	Dbauditsegment           Dbauditsegment           `yaml:"dbauditsegment"`
	Dbauditcont              Dbauditcont              `yaml:"dbauditcont"`
	Db_Nosys_In_System       Db_Nosys_In_System       `yaml:"db_nosys_in_system"`
	Dbvirscheck              Dbvirscheck              `yaml:"dbvirscheck"`
	Dbscnhealthcheck         Dbscnhealthcheck         `yaml:"dbscnhealthcheck"`
	Dbdbapriv                Dbdbapriv                `yaml:"dbdbapriv"`
	Dbsysdba                 Dbsysdba                 `yaml:"dbsysdba"`
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
	Open_files_ne        int `yaml:"open_files_ne"`
	Max_user_rocesses_ne int `yaml:"max_user_rocesses_ne"`
}

type Filesystem struct {
	Disk_ge1  string `yaml:"disk_ge1"`
	Disk_ge2  string `yaml:"disk_ge2"`
	Inode_ge1 string `yaml:"inode_ge1"`
	Inode_ge2 string `yaml:"inode_ge2"`
}

type Cpustat struct {
	Idle_le1 int `yaml:"idle_le1"`
	Idle_le2 int `yaml:"idle_le2"`
	Swap_ge1 int `yaml:"swap_ge1"`
	Swap_ge2 int `yaml:"swap_ge2"`
}
type Memstat struct {
	Available_le1 int `yaml:"available_le1"`
	Available_le2 int `yaml:"available_le2"`
}

type Iostat struct {
	Diskutil_ge1 float64 `yaml:"diskutil_ge1"`
	Diskutil_ge2 float64 `yaml:"diskutil_ge2"`
}

type Thpstat struct {
	Anpages_gt int `yaml:"anpages_gt"`
}

type Numa struct {
	Flg1 string `yaml:"flg1"`
	Flg2 string `yaml:"flg2"`
}

/////*** Lv3 os End***/////

// /////解析rule.yaml中的数据库部份规则
// ///*** Lv3 db Start***/////
type DbTbsusage struct {
	Tbsutil_ge1  float64 `yaml:"tbsutil_ge1"`
	Tbsutil_ge2  float64 `yaml:"tbsutil_ge2"`
	Freesize_le1 float64 `yaml:"freesize_le1"`
	Freesize_le2 float64 `yaml:"freesize_le2"`
}

type Dbdatafile struct {
	Status string `yaml:"status"`
}

type Dbcontrolfile struct {
	Cnt_le1 int `yaml:"cnt_le1"`
}
type Dbredocheck struct {
	Rdf_size_lt1 float64 `yaml:"rdf_size_lt1"`
	// Rdf_status   string `yaml:"rdf_status"`
	Rdf_status_list []string `yaml:"rdf_status_list,flow"` //返回字符串数组, flow为固定词
}

type Dbredoswitch struct {
	Sw_cnt_ge1 int `yaml:"sw_cnt_ge1"`
}

type Dbresource struct {
	Res_use_ge1 int `yaml:"res_use_ge1"`
}

type Loadprofile struct {
	Redosize_ge1 float64 `yaml:"redosize_ge1"`
	Logons_ge1   float64 `yaml:"logons_ge1"`
}

type Instefficiency struct {
	Buffer_hit  float64 `yaml:"buffer_hit"`
	Library_hit float64 `yaml:"library_hit"`
	Soft_parse  float64 `yaml:"soft_parse"`
}

type Dblsnrinfo struct {
	Log_size int `yaml:"log_size"`
}

type Dbtableparallel struct {
	Result string `yaml:"result"`
}

type Dbindexparallel struct {
	Result string `yaml:"result"`
}

type Dbinvalidindex struct {
	Result string `yaml:"result"`
}

type Dbsequence struct {
	Result string `yaml:"result"`
}

type Dbrecoverydest struct {
	Result string `yaml:"result"`
}

type Dbflashrecoveryuseage struct {
	Useage1 float64 `yaml:"useage1"`
	Useage2 float64 `yaml:"useage2"`
}

type Dberrlog struct {
	ResultB string `yaml:"resultB"`
}

type Dbproductuserfailedlogin struct {
	Result string `yaml:"result"`
}

type Dbdglagcheck struct {
	ResultB int `yaml:"resultB"`
}

type Dbdgerrcheck struct {
	ResultG string `yaml:"resultG"`
}

type Dbrmancheck struct {
	ResultB string `yaml:"resultB"`
	ResultR string `yaml:"resultR"`
}

type Dbauditsegment struct {
	ResultG string `yaml:"resultG"`
}

type Dbauditcont struct {
	ResultG int `yaml:"resultG"`
}

type Db_Nosys_In_System struct {
	ResultB string `yaml:"resultB"`
}

type Dbvirscheck struct {
	ResultR string `yaml:"resultR"`
}

type Dbscnhealthcheck struct {
	Resulta string `yaml:"resulta"`
	Resultb string `yaml:"resultb"`
	Resultc string `yaml:"resultc"`
}

type Dbdbapriv struct {
	ResultG string `yaml:"resultG"`
}

type Dbsysdba struct {
	ResultB string `yaml:"resultB"`
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
