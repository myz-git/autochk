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
	Osparameter      Osparameter      `yaml:"osparameter"`
	Ulimit           Ulimit           `yaml:"ulimit"`
	Filesystem       Filesystem       `yaml:"filesystem"`
	Inodeusage       Inodeusage       `yaml:"inodeusage"`
	Cpustat          Cpustat          `yaml:"cpustat"`
	Memstat          Memstat          `yaml:"memstat"`
	Iostat           Iostat           `yaml:"iostat"`
	Thpstat          Thpstat          `yaml:"thpstat"`
	Hugepage         Hugepage         `yaml:"hugepage"`
	Numa             Numa             `yaml:"numa"`
	Ntp              Ntp              `yaml:"ntp"`
	Tmzone           Tmzone           `yaml:"tmzone"`
	Selinux          Selinux          `yaml:"selinux"`
	Firewall         Firewall         `yaml:"firewall"`
	Nsswitch         Nsswitch         `yaml:"nsswitch"`
	Lo_mtu           Lo_mtu           `yaml:"lo_mtu"`
	Machine_platform Machine_platform `yaml:"machine_platform"`
	CPU_PERF_MODE    CPU_PERF_MODE    `yaml:"cpu_perf_mode"`
	NOZEROCONF       NOZEROCONF       `yaml:"nozeroconf"`
	RPM_PACKAGES     RPM_PACKAGES     `yaml:"rpm_packages"`
	Oslog            Oslog            `yaml:"oslog"`
}

// //*** Lv2***////
type Dbrule struct {
	// 数据库分析
	Dbstatus      Dbstatus      `yaml:"dbstatus"`
	Logmode       Logmode       `yaml:"logmode"`
	Db_lang       Db_lang       `yaml:"db_lang"`
	Dbtbsusage    Dbtbsusage    `yaml:"dbtbsusage"`
	Dbcursize     Dbcursize     `yaml:"dbcursize"`
	Dbf_size      Dbf_size      `yaml:"dbf_size"`
	Dbf_cnt       Dbf_cnt       `yaml:"dbf_cnt"`
	Dbf_stat      Dbf_stat      `yaml:"dbf_stat"`
	Tmpfile_size  Tmpfile_size  `yaml:"tmpfile_size"`
	Dbcontrolfile Dbcontrolfile `yaml:"dbcontrolfile"`
	User_info     User_info     `yaml:"user_info"`
	User_size     User_size     `yaml:"user_size"`
	Tab_info      Tab_info      `yaml:"tab_info"`
	Tab_parallel  Tab_parallel  `yaml:"tab_parallel"`
	Inx_parallel  Inx_parallel  `yaml:"inx_parallel"`
	Invalid_obj   Invalid_obj   `yaml:"invalid_obj"`
	Invalid_inx   Invalid_inx   `yaml:"invalid_inx"`
	Dbsequence    Dbsequence    `yaml:"dbsequence"`
	Db_seq_usage  Db_seq_usage  `yaml:"db_seq_usage"`

	// 数据库性能分析
	Db_shp_size    Db_shp_size    `yaml:"db_shp_size"`
	Db_shp_pct     Db_shp_pct     `yaml:"db_shp_pct"`
	Db_4031check   Db_4031check   `yaml:"db_4031check"`
	Dbresource     Dbresource     `yaml:"dbresource"`
	Loadprofile    Loadprofile    `yaml:"loadprofile"`
	Instefficiency Instefficiency `yaml:"instefficiency"`
	Topevent       Topevent       `yaml:"topevent"`
	Topsql_by_ela  Topsql_by_ela  `yaml:"topsql_by_ela"`

	// 数据库实例分析
	Cursor_share_mem  Cursor_share_mem  `yaml:"cursor_share_mem"`
	Dbredocheck       Dbredocheck       `yaml:"dbredocheck"`
	Dbredoswitch      Dbredoswitch      `yaml:"dbredoswitch"`
	Dbparameter       Dbparameter       `yaml:"dbparameter"`
	Db_parameter_file Db_parameter_file `yaml:"db_parameter_file"`
	Recovery_usage    Recovery_usage    `yaml:"recovery_usage"`
	Recovery_detail   Recovery_detail   `yaml:"recovery_detail"`

	// 数据库安全检查
	Db_expir_user      Db_expir_user      `yaml:"db_expir_user"`
	Db_password_verif  Db_password_verif  `yaml:"db_password_verif"`
	Userfailedlogin    Userfailedlogin    `yaml:"userfailedlogin"`
	Dbdbapriv          Dbdbapriv          `yaml:"dbdbapriv"`
	Dbsysdba           Dbsysdba           `yaml:"dbsysdba"`
	Dbauditsegment     Dbauditsegment     `yaml:"dbauditsegment"`
	Dbauditcont        Dbauditcont        `yaml:"dbauditcont"`
	Db_Nosys_In_System Db_Nosys_In_System `yaml:"db_nosys_in_system"`
	Dbvirscheck        Dbvirscheck        `yaml:"dbvirscheck"`
	Dbrmancheck        Dbrmancheck        `yaml:"dbrmancheck"`
	Dbscnhealthcheck   Dbscnhealthcheck   `yaml:"dbscnhealthcheck"`

	//软件使用
	Dboption   Dboption   `yaml:"dboption"`
	Dbfeatures Dbfeatures `yaml:"dbfeatures"`
	Dbpsu      Dbpsu      `yaml:"dbpsu"`
	Dbpatch    Dbpatch    `yaml:"dbpatch"`

	// 日志、集群、DataGuard、备份及杂项分析
	Dberrlog     Dberrlog     `yaml:"dberrlog"`
	Dbdglagcheck Dbdglagcheck `yaml:"dbdglagcheck"`
	Dbdgerrcheck Dbdgerrcheck `yaml:"dbdgerrcheck"`
	Dblsnrinfo   Dblsnrinfo   `yaml:"dblsnrinfo"`

	// 集群检查
	Crs_stat      Crs_stat      `yaml:"crs_stat"`
	Crs_stat2     Crs_stat2     `yaml:"crs_stat2"`
	Ocr_info      Ocr_info      `yaml:"ocr_info"`
	Ocr_bak_check Ocr_bak_check `yaml:"ocr_bak_check"`
	Asm_usage     Asm_usage     `yaml:"asm_usage"`
	Asm_offset    Asm_offset    `yaml:"asm_offset"`
}

// ///*** Lv3 OS Start***/////
type Osparameter struct {
	Nm                        string   `yaml:"nm"`
	Title                     string   `yaml:"title"`
	Desc                      string   `yaml:"desc"`
	File_max                  int      `yaml:"file_max"`
	Aio_max_nr                int      `yaml:"aio_max_nr"`
	Shmmni                    int      `yaml:"shmmni"`
	Shmmax                    int      `yaml:"shmmax"`
	Shmall                    int      `yaml:"shmall"`
	Sem                       []string `yaml:"sem,flow"`
	Panic_on_oops             int      `yaml:"panic_on_oops"`
	Randomize_va_space        int      `yaml:"randomize_va_space"`
	Numa_balancing            int      `yaml:"numa_balancing"`
	Min_free_kbytes           int      `yaml:"min_free_kbytes"`
	Rp_filter_all             int      `yaml:"rp_filter_all"`
	Rp_filter_default         int      `yaml:"rp_filter_default"`
	Ip_local_port_range       []string `yaml:"ip_local_port_range,flow"`
	Ipfrag_high_thresh        int      `yaml:"ipfrag_high_thresh"`
	Ipfrag_low_thresh         int      `yaml:"ipfrag_low_thresh"`
	Tcp_keepalive_time        int      `yaml:"tcp_keepalive_time"`
	Rmem_default              int      `yaml:"rmem_default"`
	Rmem_max                  int      `yaml:"rmem_max"`
	Wmem_default              int      `yaml:"wmem_default"`
	Wmem_max                  int      `yaml:"wmem_max"`
	Swappiness                int      `yaml:"swappiness"`
	Dirty_ratio               int      `yaml:"dirty_ratio"`
	Dirty_background_ratio    int      `yaml:"dirty_background_ratio"`
	Dirty_expire_centisecs    int      `yaml:"dirty_expire_centisecs"`
	Dirty_writeback_centisecs int      `yaml:"dirty_writeback_centisecs"`
	Disable_ism_large_pages   []string `yaml:"disable_ism_large_pages,flow"` //返回字符串数组, flow为固定词
}

type Ulimit struct {
	Nm                string `yaml:"nm"`
	Title             string `yaml:"title"`
	Desc              string `yaml:"desc"`
	Memlock           int    `yaml:"memlock"`
	Open_files        int    `yaml:"open_files"`
	Max_user_rocesses int    `yaml:"max_user_rocesses"`
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

type Hugepage struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
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

type Tmzone struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Selinux struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Firewall struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Nsswitch struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Lo_mtu struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Machine_platform struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type CPU_PERF_MODE struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type NOZEROCONF struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type RPM_PACKAGES struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Oslog struct {
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
	Status string `yaml:"status"`
}

type Logmode struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Status string `yaml:"status"`
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

type Dbcursize struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbf_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbf_cnt struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbf_stat struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Status string `yaml:"status"`
}

type Tmpfile_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbcontrolfile struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Cnt_le1 int    `yaml:"cnt_le1"`
}

type User_info struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type User_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Tab_info struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Tab_parallel struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Inx_parallel struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Invalid_obj struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Invalid_inx struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Dbredocheck struct {
	Nm       string  `yaml:"nm"`
	Title    string  `yaml:"title"`
	Desc     string  `yaml:"desc"`
	Rdf_size float64 `yaml:"rdf_size"`
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
	Nm       string  `yaml:"nm"`
	Title    string  `yaml:"title"`
	Desc     string  `yaml:"desc"`
	Redosize float64 `yaml:"redosize"`
	Logon    float64 `yaml:"logon"`
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

type Topsql_by_ela struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Cursor_share_mem struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dblsnrinfo struct {
	Nm       string `yaml:"nm"`
	Title    string `yaml:"title"`
	Desc     string `yaml:"desc"`
	Log_size int    `yaml:"log_size"`
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

type Recovery_usage struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result []int  `yaml:"result"`
}

type Recovery_detail struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
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

type Db_password_verif struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
}

type Userfailedlogin struct {
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
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
}

type Dbauditcont struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
}

type Db_Nosys_In_System struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
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
	Nm     string  `yaml:"nm"`
	Title  string  `yaml:"title"`
	Desc   string  `yaml:"desc"`
	Result float64 `yaml:"result"`
}

type Db_4031check struct {
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

type Dboption struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Dbfeatures struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Db_lang struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Crs_stat struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Crs_stat2 struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Ocr_info struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Ocr_bak_check struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Asm_usage struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Asm_offset struct {
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
	configFile, err = os.ReadFile("./local/rule.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
}
