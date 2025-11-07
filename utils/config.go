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
	Hostname         Hostname         `yaml:"hostname"`
	Ipaddr           Ipaddr           `yaml:"ipaddr"`
	Os               Os               `yaml:"os"`
	Relver           Relver           `yaml:"relver"`
	Cpu_model        Cpu_model        `yaml:"cpu_model"`
	Cpucount         Cpucount         `yaml:"cpucount"`
	Cpumhz           Cpumhz           `yaml:"cpumhz"`
	Memtotal         Memtotal         `yaml:"memtotal"`
	Swaptotal        Swaptotal        `yaml:"swaptotal"`
	Osparam_fs       Osparam_fs       `yaml:"osparam_fs"`
	Osparam_ker      Osparam_ker      `yaml:"osparam_ker"`
	Osparam_net      Osparam_net      `yaml:"osparam_net"`
	Osparam_vm       Osparam_vm       `yaml:"osparam_vm"`
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
	Dbstatus   Dbstatus   `yaml:"dbstatus"`
	Logmode    Logmode    `yaml:"logmode"`
	Db_lang    Db_lang    `yaml:"db_lang"`
	Dbtbsusage Dbtbsusage `yaml:"dbtbsusage"`
	Dbcursize  Dbcursize  `yaml:"dbcursize"`
	Dbf_size   Dbf_size   `yaml:"dbf_size"`
	Dbf_cnt    Dbf_cnt    `yaml:"dbf_cnt"`

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
	Dbparam_b         Dbparam_b         `yaml:"dbparam_b"`
	Dbparam_d         Dbparam_d         `yaml:"dbparam_d"`
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
type Hostname struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Ipaddr struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Os struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Relver struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Cpu_model struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Cpucount struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Cpumhz struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Memtotal struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Swaptotal struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Osparam_fs struct {
	Nm         string `yaml:"nm"`
	Title      string `yaml:"title"`
	Desc       string `yaml:"desc"`
	File_max   int    `yaml:"file_max"`
	Aio_max_nr int    `yaml:"aio_max_nr"`
	Level      string `yaml:"level"`
}

type Osparam_ker struct {
	Nm                 string   `yaml:"nm"`
	Title              string   `yaml:"title"`
	Desc               string   `yaml:"desc"`
	Shmmni             int      `yaml:"shmmni"`
	Shmmax             int      `yaml:"shmmax"`
	Shmall             int      `yaml:"shmall"`
	Sem                []string `yaml:"sem,flow"`
	Panic_on_oops      int      `yaml:"panic_on_oops"`
	Randomize_va_space int      `yaml:"randomize_va_space"`
	Numa_balancing     int      `yaml:"numa_balancing"`
	Level              string   `yaml:"level"`
}

type Osparam_net struct {
	Nm                  string   `yaml:"nm"`
	Title               string   `yaml:"title"`
	Desc                string   `yaml:"desc"`
	Rp_filter_all       int      `yaml:"rp_filter_all"`
	Rp_filter_default   int      `yaml:"rp_filter_default"`
	Ip_local_port_range []string `yaml:"ip_local_port_range,flow"`
	Ipfrag_high_thresh  int      `yaml:"ipfrag_high_thresh"`
	Ipfrag_low_thresh   int      `yaml:"ipfrag_low_thresh"`
	Tcp_keepalive_time  int      `yaml:"tcp_keepalive_time"`
	Rmem_default        int      `yaml:"rmem_default"`
	Rmem_max            int      `yaml:"rmem_max"`
	Wmem_default        int      `yaml:"wmem_default"`
	Wmem_max            int      `yaml:"wmem_max"`
	Level               string   `yaml:"level"`
}

type Osparam_vm struct {
	Nm                        string   `yaml:"nm"`
	Title                     string   `yaml:"title"`
	Desc                      string   `yaml:"desc"`
	Swappiness                int      `yaml:"swappiness"`
	Min_free_kbytes           int      `yaml:"min_free_kbytes"`
	Dirty_ratio               int      `yaml:"dirty_ratio"`
	Dirty_background_ratio    int      `yaml:"dirty_background_ratio"`
	Dirty_expire_centisecs    int      `yaml:"dirty_expire_centisecs"`
	Dirty_writeback_centisecs int      `yaml:"dirty_writeback_centisecs"`
	Disable_ism_large_pages   []string `yaml:"disable_ism_large_pages,flow"` //返回字符串数组, flow为固定词
	Level                     string   `yaml:"level"`
}

type Ulimit struct {
	Nm                string `yaml:"nm"`
	Title             string `yaml:"title"`
	Desc              string `yaml:"desc"`
	Memlock           int    `yaml:"memlock"`
	Open_files        int    `yaml:"open_files"`
	Max_user_rocesses int    `yaml:"max_user_rocesses"`
	Level             string `yaml:"level"`
}

type Filesystem struct {
	Nm       string `yaml:"nm"`
	Title    string `yaml:"title"`
	Desc     string `yaml:"desc"`
	Disk_ge1 string `yaml:"disk_ge1"`
	Disk_ge2 string `yaml:"disk_ge2"`
	Level    string `yaml:"level"`
}

type Inodeusage struct {
	Nm        string `yaml:"nm"`
	Title     string `yaml:"title"`
	Desc      string `yaml:"desc"`
	Inode_ge1 string `yaml:"inode_ge1"`
	Inode_ge2 string `yaml:"inode_ge2"`
	Level     string `yaml:"level"`
}

type Cpustat struct {
	Nm       string `yaml:"nm"`
	Title    string `yaml:"title"`
	Desc     string `yaml:"desc"`
	Idle_le1 int    `yaml:"idle_le1"`
	Idle_le2 int    `yaml:"idle_le2"`
	Swap_ge1 int    `yaml:"swap_ge1"`
	Swap_ge2 int    `yaml:"swap_ge2"`
	Level    string `yaml:"level"`
}
type Memstat struct {
	Nm            string `yaml:"nm"`
	Title         string `yaml:"title"`
	Desc          string `yaml:"desc"`
	Available_le1 int    `yaml:"available_le1"`
	Available_le2 int    `yaml:"available_le2"`
	Level         string `yaml:"level"`
}

type Iostat struct {
	Nm           string  `yaml:"nm"`
	Title        string  `yaml:"title"`
	Desc         string  `yaml:"desc"`
	Diskutil_ge1 float64 `yaml:"diskutil_ge1"`
	Diskutil_ge2 float64 `yaml:"diskutil_ge2"`
	Level        string  `yaml:"level"`
}

type Thpstat struct {
	Nm         string `yaml:"nm"`
	Title      string `yaml:"title"`
	Desc       string `yaml:"desc"`
	Anpages_gt int    `yaml:"anpages_gt"`
	Level      string `yaml:"level"`
}

type Hugepage struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Numa struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Flg1  string `yaml:"flg1"`
	Flg2  string `yaml:"flg2"`
	Level string `yaml:"level"`
}

type Ntp struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Tmzone struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Selinux struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Firewall struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Nsswitch struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Lo_mtu struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Machine_platform struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type CPU_PERF_MODE struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type NOZEROCONF struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type RPM_PACKAGES struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Oslog struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

/////*** Lv3 os End***/////

// /////解析rule.yaml中的数据库部份规则
// ///*** Lv3 db Start***/////

type Dbstatus struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Status string `yaml:"status"`
	Level  string `yaml:"level"`
}

type Logmode struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Status string `yaml:"status"`
	Level  string `yaml:"level"`
}

type Dbtbsusage struct {
	Nm           string  `yaml:"nm"`
	Title        string  `yaml:"title"`
	Desc         string  `yaml:"desc"`
	Tbsutil_ge1  float64 `yaml:"tbsutil_ge1"`
	Tbsutil_ge2  float64 `yaml:"tbsutil_ge2"`
	Freesize_le1 float64 `yaml:"freesize_le1"`
	Freesize_le2 float64 `yaml:"freesize_le2"`
	Level        string  `yaml:"level"`
}

type Dbcursize struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbf_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbf_cnt struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbf_stat struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Status string `yaml:"status"`
	Result int    `yaml:"result"`
	Level  string `yaml:"level"`
}

type Tmpfile_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbcontrolfile struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Cnt_le1 int    `yaml:"cnt_le1"`
	Level   string `yaml:"level"`
}

type User_info struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type User_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Tab_info struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Tab_parallel struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
	Level  string `yaml:"level"`
}

type Inx_parallel struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
	Level  string `yaml:"level"`
}

type Invalid_obj struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
	Level  string `yaml:"level"`
}

type Invalid_inx struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
	Level  string `yaml:"level"`
}

type Dbredocheck struct {
	Nm       string  `yaml:"nm"`
	Title    string  `yaml:"title"`
	Desc     string  `yaml:"desc"`
	Rdf_size float64 `yaml:"rdf_size"`
	// Rdf_status   string `yaml:"rdf_status"`
	Rdf_status_list []string `yaml:"rdf_status_list,flow"` //返回字符串数组, flow为固定词
	Level           string   `yaml:"level"`
}

type Dbredoswitch struct {
	Nm         string `yaml:"nm"`
	Title      string `yaml:"title"`
	Desc       string `yaml:"desc"`
	Sw_cnt_ge1 int    `yaml:"sw_cnt_ge1"`
	Level      string `yaml:"level"`
}

type Dbresource struct {
	Nm          string `yaml:"nm"`
	Title       string `yaml:"title"`
	Desc        string `yaml:"desc"`
	Res_use_ge1 int    `yaml:"res_use_ge1"`
	Level       string `yaml:"level"`
}

type Loadprofile struct {
	Nm       string  `yaml:"nm"`
	Title    string  `yaml:"title"`
	Desc     string  `yaml:"desc"`
	Redosize float64 `yaml:"redosize"`
	Logon    float64 `yaml:"logon"`
	Level    string  `yaml:"level"`
}

type Instefficiency struct {
	Nm          string  `yaml:"nm"`
	Title       string  `yaml:"title"`
	Desc        string  `yaml:"desc"`
	Buffer_hit  float64 `yaml:"buffer_hit"`
	Library_hit float64 `yaml:"library_hit"`
	Soft_parse  float64 `yaml:"soft_parse"`
	Level       string  `yaml:"level"`
}

type Topevent struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Topsql_by_ela struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Cursor_share_mem struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dblsnrinfo struct {
	Nm       string `yaml:"nm"`
	Title    string `yaml:"title"`
	Desc     string `yaml:"desc"`
	Log_size int    `yaml:"log_size"`
	Level    string `yaml:"level"`
}

type Dbsequence struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Db_seq_usage struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Recovery_usage struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result []int  `yaml:"result"`
	Level  string `yaml:"level"`
}

type Recovery_detail struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Dberrlog struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB string `yaml:"resultB"`
	Level   string `yaml:"level"`
}

type Db_expir_user struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Db_password_verif struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Userfailedlogin struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Dbdbapriv struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultG string `yaml:"resultG"`
	Level   string `yaml:"level"`
}

type Dbsysdba struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB string `yaml:"resultB"`
	Level   string `yaml:"level"`
}

type Dbdglagcheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB int    `yaml:"resultB"`
	Level   string `yaml:"level"`
}

type Dbdgerrcheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultG string `yaml:"resultG"`
	Level   string `yaml:"level"`
}

type Dbauditsegment struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
	Level  string `yaml:"level"`
}

type Dbauditcont struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result int    `yaml:"result"`
	Level  string `yaml:"level"`
}

type Db_Nosys_In_System struct {
	Nm     string `yaml:"nm"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Result string `yaml:"result"`
	Level  string `yaml:"level"`
}

type Dbvirscheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultR string `yaml:"resultR"`
	Level   string `yaml:"level"`
}

type Dbscnhealthcheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	Resulta string `yaml:"resulta"`
	Resultb string `yaml:"resultb"`
	Resultc string `yaml:"resultc"`
	Level   string `yaml:"level"`
}

type Dbparameter struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Db_parameter_file struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Db_shp_size struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Db_shp_pct struct {
	Nm     string  `yaml:"nm"`
	Title  string  `yaml:"title"`
	Desc   string  `yaml:"desc"`
	Result float64 `yaml:"result"`
	Level  string  `yaml:"level"`
}

type Db_4031check struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbrmancheck struct {
	Nm      string `yaml:"nm"`
	Title   string `yaml:"title"`
	Desc    string `yaml:"desc"`
	ResultB string `yaml:"resultB"`
	ResultR string `yaml:"resultR"`
	Level   string `yaml:"level"`
}

type Dbpsu struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbpatch struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbparam_b struct {
	Nm                      string `yaml:"nm"`
	Title                   string `yaml:"title"`
	Desc                    string `yaml:"desc"`
	Level                   string `yaml:"level"`
	Db_cache_size           int    `yaml:"db_cache_size"`
	Db_files                int    `yaml:"db_files"`
	Memory_max_target       int    `yaml:"memory_max_target"`
	Memory_target           int    `yaml:"memory_target"`
	Open_cursors            int    `yaml:"open_cursors"`
	Open_links              int    `yaml:"open_links"`
	Open_links_per_instance int    `yaml:"open_links_per_instance"`
	Pga_aggregate_target    int    `yaml:"pga_aggregate_target"`
	Processes               int    `yaml:"processes"`
	Session_cached_cursors  int    `yaml:"session_cached_cursors"`
	Sga_max_size            int    `yaml:"sga_max_size"`
	Sga_target              int    `yaml:"sga_target"`
	Shared_pool_size        int    `yaml:"shared_pool_size"`
	Streams_pool_size       int    `yaml:"streams_pool_size"`
	Undo_retention          int    `yaml:"undo_retention"`
}

type Dbparam_d struct {
	Nm                                      string `yaml:"nm"`
	Title                                   string `yaml:"title"`
	Desc                                    string `yaml:"desc"`
	Level                                   string `yaml:"level"`
	U_and_pruning_enabled                   string `yaml:"_and_pruning_enabled"`
	U_ash_size                              int    `yaml:"_ash_size"`
	U_bloom_filter_enabled                  string `yaml:"_bloom_filter_enabled"`
	U_bloom_pruning_enabled                 string `yaml:"_bloom_pruning_enabled"`
	U_cleanup_rollback_entries              int    `yaml:"_cleanup_rollback_entries"`
	U_cursor_obsolete_threshold             int    `yaml:"_cursor_obsolete_threshold"`
	U_datafile_write_errors_crash_instance  string `yaml:"_datafile_write_errors_crash_instance"`
	U_disable_last_successful_login_time    string `yaml:"_disable_last_successful_login_time"`
	U_drop_stat_segment                     int    `yaml:"_drop_stat_segment"`
	U_lm_tickets                            int    `yaml:"_lm_tickets"`
	U_max_spacebg_slaves                    int    `yaml:"_max_spacebg_slaves"`
	U_optimizer_adaptive_cursor_sharing     string `yaml:"_optimizer_adaptive_cursor_sharing"`
	U_optimizer_extended_cursor_sharing     string `yaml:"_optimizer_extended_cursor_sharing"`
	U_optimizer_extended_cursor_sharing_rel string `yaml:"_optimizer_extended_cursor_sharing_rel"`
	U_optimizer_null_accepting_semijoin     string `yaml:"_optimizer_null_accepting_semijoin"`
	U_optimizer_outer_to_anti_enabled       string `yaml:"_optimizer_outer_to_anti_enabled"`
	U_optimizer_partial_join_eval           string `yaml:"_optimizer_partial_join_eval"`
	U_optimizer_reduce_groupby_key          string `yaml:"_optimizer_reduce_groupby_key"`
	U_optimizer_use_feedback                string `yaml:"_optimizer_use_feedback"`
	U_optimizer_gather_feedback             string `yaml:"_optimizer_gather_feedback"`
	U_rowsets_enabled                       string `yaml:"_rowsets_enabled"`
	U_report_capture_cycle_time             int    `yaml:"_report_capture_cycle_time"`
	U_securefiles_concurrency_estimate      int    `yaml:"_securefiles_concurrency_estimate"`
	U_shared_pool_reserved_pct              int    `yaml:"_shared_pool_reserved_pct"`
	U_sys_logon_delay                       int    `yaml:"_sys_logon_delay"`
	U_use_adaptive_log_file_sync            string `yaml:"_use_adaptive_log_file_sync"`
	U_undo_autotune                         string `yaml:"_undo_autotune"`
	U_use_single_log_writer                 string `yaml:"_use_single_log_writer"`
	Client_statistics_level                 string `yaml:"client_statistics_level"`
	Control_file_record_keep_time           int    `yaml:"control_file_record_keep_time"`
	Deferred_segment_creation               string `yaml:"deferred_segment_creation"`
	Enable_ddl_logging                      string `yaml:"enable_ddl_logging"`
	Fast_start_mttr_target                  int    `yaml:"fast_start_mttr_target"`
	Max_dump_file_size                      int    `yaml:"max_dump_file_size"`
	Parallel_execution_message_size         int    `yaml:"parallel_execution_message_size"`
	Parallel_force_local                    string `yaml:"parallel_force_local"`
	Parallel_max_servers                    int    `yaml:"parallel_max_servers"`
}

type Dboption struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Dbfeatures struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Db_lang struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Crs_stat struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Crs_stat2 struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Ocr_info struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Ocr_bak_check struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Asm_usage struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
}

type Asm_offset struct {
	Nm    string `yaml:"nm"`
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
	Level string `yaml:"level"`
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
