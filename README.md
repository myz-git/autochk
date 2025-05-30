# autochk
健康检查自动分析系统

# 主要功能介绍: 
根据预设的检查规则对收集到的系统及数据库信息(XML格式文件)进行读取和分析,自动判断是否满足规则,然后最终生成检查报告;

# 详细功能设计:
1. 给定检查规则rule.yaml文件,内容为OS及DB的检查项名称及检查值健康判断标准;
2. autochk对给定的xml文件(内容为对LINUX及数据库的信息收集结果) 进行读取,在读取过程中会根据检查规则对每一个XML中的检查项进行健康判断,然后将检查项,检查结果,健康判断等保存到xlsx表中;
3. 最后将xlsx表根据WORD文件模板生成健康检查报告


# 项目架构:
|--main.go	程序总入口,调度readxml读取xml,调度anadata分析xml内容,调度toxls写入xlsx文件;
|
|-structs/structs.go 定义了InfoSht,OsSht,DbSht三个structs,分别对应"展现结构体","系统结构体","数据库结构体";
|
|-readxml/readxml.go 读取给定的xml文件,遍历tag标签,提取出每个检查项目及内容,分类保存到infoshtp,osshtp,dbshtp三个struct中;
|
|-anadata/ana.go	建立分析函数,对struct中的检查项目及内容,对照rule.yaml中定义的检查项及健康判定规则进行数据分析,将健康判断分析结果再保存到structs中;
|
|-toxls/xlsx.go 将经过ana.go分析后的InfoSht,OsSht,DbSht三个structs内容,格式化写入到xlsx文件中;
|
|-utils/config.go 解析yaml文件获得检查项及健康判定规则;
|
|-rule.yaml  健康检查规则文件;
|
|-20230322_dcs0tdb6_apc2db.ALL.xml  xml样例文件;




# 程序实现:
##  使用 etree 解析复杂结构的 xml 文件
参考如下
https://godoc.org/github.com/beevik/etree
https://pkg.go.dev/github.com/beevik/etree?tab=doc
https://github.com/beevik/etree

## 填加指标步骤:
1. 在stucts.go中填加type
2. 在xlsx.NewXlsx 写入初始化值, 并相应修改sheet下的style单元格生效范围
3. 在xlsx.PutSht_INFO(或PutSht_OS,PutSht_DB)中 switch rowCell 添加相应case 中文要和newxlsx中一致;
4. 在readXml.go tag0或tag1相应位置添加解析
5. 如果是检查指标,在rule.yaml中设定检查规则 (le ge ne lt gt 分别对应 <=  >= !=  <  >  1,2代表级别),注意":"后要有空格如  nm: value
6. 如果是检查指标,在confg中添加检查规则名称, 要和yaml完全一致, 注意数据类型
7. 在ana中添加指标分析函数Ana_xxx或者格式化输出函数Fmt_xxx

## 程序说明
### main.go
命令行参数处理：通过 flag 包处理 -s 参数，决定是否以单文件模式处理XML文件。
清理旧文件：执行 ClearFile 函数来删除已完成的 .xlsx 文件，以准备新的输出。
文件获取：使用 GetXMLS 函数遍历目录，获取指定类型的XML文件列表。
文件处理循环：
对每个 .ALL.xml 文件初始化展示结构体（InfoSht）、系统结构体（OsSht）、数据库结构体（DbSht）。
调用 readxml.ReadXml 读取并填充结构体。
调用 anadata.Ana 分析结构体中的数据。
调用 toxls.Xlsx 生成Excel文件。
性能监控：记录程序运行时间，打印开始与结束日志，方便监控程序执行情况。

### structs/structs.go
structs.go定义了三个主要的结构体，用于存储系统和数据库检查的结果：

InfoSht：这个结构体用于存储数据库和操作系统的基本信息，如数据库名称、版本、角色、日志模式、大小、文件数量、表数量、语言等，以及操作系统的相关信息，如主机名、IP地址、操作系统版本、CPU核心数、总内存等。
OsSht：这个结构体包含了多个 Tpstrc 类型的字段，每个字段代表操作系统的一个特定参数或状态，如操作系统参数、文件系统、CPU状态、内存状态等。这些信息将用于评估系统的健康状态。
DbSht：类似于 OsSht，这个结构体包含了多个关于数据库的检查点，如表空间使用情况、数据文件、控制文件、重做日志检查、SQL性能、监听器信息等。每个检查点都使用 Tpstrc 结构体存储内容和告警级别。
Tpstrc：这是一个辅助结构体，用于存储具体的检查内容和相应的告警级别（红色、蓝色、绿色），以标示不同程度的健康风险。
这些结构体为系统和数据库的健康检查提供了一个清晰的数据模型，便于在程序中处理和传递检查数据。

### rule.yaml
rule.yaml文件详细定义了操作系统（OS）和数据库（DB）检查规则，为每个检查项目设置了特定的阈值和条件。这些规则是 anadata 模块使用来评估系统和数据库健康状况的依据。以下是对文件中一些关键部分的解读：

OS检查规则 (osrule):
osparameter：定义了操作系统级别的参数检查，如最大进程数、文件数、内存最小空闲等。
ulimit：设置了用户限制，如打开文件的最大数量和最大用户进程数。
filesystem：定义了磁盘和inode的使用率阈值。
cpustat：涵盖了CPU空闲率和交换空间的阈值。
memstat：关于内存可用量的阈值。
iostat：磁盘利用率的检查。
thpstat、numa：涉及透明大页面和NUMA配置的检查。
DB检查规则 (dbrule):
dbtbsusage：表空间使用情况，包括使用率和剩余空间的阈值。
dbdatafile、dbcontrolfile：关于数据文件和控制文件的状态和数量检查。
dbredocheck、dbredoswitch：涉及日志文件大小和切换频率的检查。
dbresource、loadprofile：资源使用情况和系统负载相关的检查。
instefficiency：数据库实例效率相关的检查。
其他项目：包括监听信息、并行处理、索引有效性等的检查。

### utils/config.go
config.go文件负责解析 rule.yaml 文件，并将其中的配置加载到预定义的结构体中。这里的结构体层次清晰地映射了 YAML 文件中的数据结构，从而允许程序以类型安全的方式访问这些配置数据。以下是文件中关键部分的概览：

主要结构体
RuleInfo：顶层结构体，包含操作系统和数据库规则的结构体。
Osrule 和 Dbrule：分别定义操作系统和数据库的具体规则结构体，这些结构体映射了 YAML 中相应的数据层级。
操作系统（OS）规则结构体
包括 Osparameter, Ulimit, Filesystem, Cpustat, Memstat, Iostat, Thpstat, Numa 等，每个结构体映射到 YAML 中的对应部分，用于定义具体的检查参数和阈值。

数据库（DB）规则结构体
包括 DbTbsusage, Dbdatafile, Dbcontrolfile, Dbredocheck 等，每个结构体映射到 YAML 中的对应部分，同样用于定义数据库健康检查的具体参数和阈值。

解析函数
GetRule：该函数使用 yaml.Unmarshal 方法将 rule.yaml 的内容解析到 RuleInfo 结构体实例中。这允许其他部分的程序通过这个实例访问配置数据。
初始化函数
init：在程序启动时自动调用，用于从文件系统中读取 rule.yaml 文件并加载到 configFile 变量中。如果文件读取失败，程序将记录错误并终止。
注意事项
使用 os.ReadFile 替代了旧的 ioutil.ReadFile，这是一个改进的做法，因为 ioutil 包在 Go 1.16 后被标记为废弃。
这样的设计确保了程序的配置数据在启动时被加载和解析，而且由于采用了结构化的方式，增加新的配置规则或修改现有规则都变得相对简单

### readxml/readxml.go
readxml.go文件包含了 ReadXml 函数，这个函数用于读取 XML 文件并解析其内容，填充到相应的结构体中。这里使用了 etree 包来处理 XML，这是一个在 Go 中处理 XML 数据的常用库。
以下是对代码的详细分析：

功能概述
初始化和读取XML文件：
使用 etree.NewDocument() 创建一个新的 XML 文档对象。
使用 doc.ReadFromFile(path) 从指定路径读取 XML 文件。
解析XML结构：
首先获取根元素 EACHK。
在 TAG0 和 TAG1 中遍历 XML 标签，这些标签似乎对应不同类型的数据结构。
填充结构体：
根据标签的类型，将数据填充到 InfoSht, OsSht, 或 DbSht 结构体中。每个标签对应结构体中的一个字段。
主要代码部分
填充 InfoSht 结构体：
处理与主机信息相关的标签，如 HOSTNAME, IPADDR, OS, CORES 等。
填充 OsSht 结构体：
处理与操作系统设置相关的标签，如 OSPARAMETER, ULIMIT, FILESYSTEM 等。
填充 DbSht 结构体：
处理与数据库相关的标签，如 DBNAME, DBVER, DBROLE, LOGMODE 以及其他数据库性能和配置相关的标签。
注意事项和改进建议
错误处理：
使用 panic(err) 来处理文件读取错误。虽然这可以在出错时立即停止程序，但在生产环境中通常建议使用更温和的错误处理策略，比如返回错误到上层调用者，以允许更优雅的错误处理和恢复。
代码组织：
如果 XML 结构复杂或标签种类多，可以考虑将处理逻辑拆分成多个辅助函数，以保持 ReadXml 函数的简洁和可维护性。
性能优化：
当处理非常大的 XML 文件时，可能需要考虑性能优化，例如使用流式处理（streaming）来减少内存使用。
这个函数是系统中非常关键的一部分，因为它直接处理输入数据，并将这些数据转换成内部可以进一步处理的格式。

### anadata/ana.go
ana.go 文件中包含的代码主要负责分析和格式化来自 XML 文件的数据，以及将这些数据与 rule.yaml 文件中定义的规则进行比较，以评估系统和数据库的健康状态。这里是对您提供的代码部分的解释和分析：

功能概述
数据转换：
String2Int 函数将字符串数组转换为整数数组。
包含检查：
Contain 函数检查一个对象是否包含在给定的集合中，支持数组、切片和映射。
主分析函数：
Ana 函数是主分析函数，它调用多个辅助函数来分析和格式化结构体中的数据。
数据格式化函数
这些函数专门用于格式化 InfoSht 结构体中的数据：

Fmt_DbRole, Fmt_LogMode, Fmt_FlashBack, Fmt_DbTotalsize, Fmt_DbFilecount, Fmt_DbTblcount：
每个函数处理不同的数据字段。
通常，这些函数从多行数据中提取第三行的信息（如果存在），并在必要时添加额外的文本（如添加 "GB" 单位）。
主要分析逻辑
在 Ana 函数中，首先尝试从 utils 包中获取规则，然后调用一系列函数来分析操作系统（OS）和数据库（DB）相关的数据。
每个 Ana_ 前缀的函数都负责分析特定的系统或数据库指标，比如 Ana_Osparameter, Ana_Ulimit, Ana_DbTbs 等。
设计和实现考虑
错误处理：在主分析函数中，当获取规则失败时，只是打印错误日志。这种处理方式可能需要根据应用的需求进一步优化，例如，可以考虑将错误返回到更高的调用层。
代码复用：格式化数据的函数有类似的结构和逻辑，可能会考虑通过创建更通用的函数来减少代码重复。
性能考虑：如果处理的数据量很大，需要评估这些字符串操作和反射调用的性能影响

Ana_Osparameter 和 Ana_Ulimit 函数：
这两个函数主要针对不同的操作系统参数进行检查。使用正则表达式提取关键信息，并与预设规则比较，根据结果设置告警级别。
Ana_Osparameter
根据提供的系统类型（如 "LINUX" 或 "SOLARIS"），函数使用正则表达式匹配关键配置项，并与规则中的预设阈值比较。
如果发现配置项的值未达标准或不符合预期，会设置相应的告警级别（如 "B" 表示蓝色重要告警，"G" 表示绿色普通告警）。
使用 break Looop 从循环中退出，表示一旦发现不符合规则的项目就停止检查后续项。
Ana_Ulimit
功能与 Ana_Osparameter 类似，专注于用户权限设置（如最大打开文件数和最大用户进程数）的检查。
针对不同的操作系统，匹配和验证设置的不同参数。
同样使用 break Looop 逻辑，一旦发现不符合规则的设置，立即设置告警并停止检查。
设计和实现的考虑点
正则表达式的使用：
正则表达式是处理文本匹配的强大工具，但也需要注意其性能影响，特别是在需要频繁编译表达式的场景中。考虑预编译正则表达式以提高性能。
错误处理：
虽然函数中使用了 _ 来忽略 strconv.Atoi 的错误返回，但在生产环境中通常建议处理这些错误以避免潜在的数据问题。
代码重构：
多个函数中存在类似的逻辑（如正则表达式匹配和阈值检查），可能可以通过抽象出共通函数或使用策略模式来简化代码并提高可维护性。
日志记录：
代码中注释掉的日志记录是一个好的调试工具，但也要确保生产环境中的日志记录既足够提供必要信息，又不会过多地影响性能或产生过大的日志文件

第三部分的 ana.go 文件中继续展示了一系列的函数，用于分析数据库相关的健康指标。这些函数利用正则表达式和字符串操作来解析提供的数据，并根据 rule.yaml 中的规则设置警报级别。以下是对您提供的函数的详细分析：

Ana_DbTbs (数据库表空间使用率检查)
主要逻辑：遍历每一行数据，检查表空间使用率，如果达到预设的告警阈值，则设置相应的告警等级并中断循环。
警报逻辑：
使用率大于等于 90% 且剩余空间小于 4GB 设置为 "R"（红色警报）。
使用率大于等于 80% 且剩余空间小于 8GB 设置为 "B"（蓝色警报）。
Ana_DBF (数据库文件状态检查)
主要逻辑：检查数据库文件的状态，如果状态不符合预设值，设置为 "R"（红色警报）并退出循环。
Ana_DBCTRF (控制文件数量检查)
主要逻辑：检查控制文件的数量，如果少于预设的最小数量，设置为 "B"（蓝色警报）。
Ana_RDF (重做日志文件检查)
主要逻辑：检查重做日志文件的状态和大小，如果不符合预设状态列表或大小小于预设值，则设置警报。
Ana_RDSW (重做日志切换频率检查)
主要逻辑：分析重做日志切换的频率，如果超过预设阈值，则设置为 "B"（蓝色警报）。
Ana_RESOURCE (资源使用率检查)
主要逻辑：对比实际资源使用率与阈值，如果超过预设的比例，设置为 "R"（红色警报）。
Ana_LOADPROFILE (负载分析)
主要逻辑：分析日志大小和登录次数，如果超过预设阈值，根据情况设置 "G"（绿色警报）或 "B"（蓝色警报）。
Ana_INSTEFFICIENCY (实例效率分析)
主要逻辑：分析缓冲区命中率、库缓存命中率和软解析比率，如果低于预设阈值，设置为 "G"（绿色警报）。
设计和实现的考虑点
正则表达式和字符串处理：
这些函数广泛使用正则表达式来解析和验证数据，这是有效的，但需要确保表达式的准确性和性能。
代码优化：
很多函数中的逻辑相似，可以考虑抽象出通用的处理逻辑或创建辅助函数来简化代码和提高重用性。
错误处理：
大多数转换函数（如 strconv.ParseFloat）的错误返回被忽略或仅在错误时调用 panic。建议处理这些错误，可能通过返回错误给调用者，而不是在底层直接导致程序崩溃。

### toxls/xlsx.go
toxls/xlsx.go 文件中，包含了从结构体中读取数据并将其保存到 Excel 文件的逻辑。文件使用了 excelize/v2 包，这是一个流行的 Go 语言库，用于读写 Excel 文件。以下是对该文件关键功能的概括和分析：

主要功能和代码流程
文件和工作表创建：
根据传入的参数（sglf 单文件标志）决定是为每个 XML 文件创建一个独立的 Excel 文件还是将所有数据保存到一个统一的 Excel 文件中。
在第一次调用时创建新的 Excel 文件。
数据写入逻辑：
使用 PutSht_INFO，PutSht_OS，和 PutSht_DB 函数将不同来源的数据写入不同的工作表（INFO, OS, DB）。
每个函数处理对应的结构体数据，设置特定的单元格内容，以及根据警报级别设置单元格样式。
样式和格式设置：
为不同的警报级别定义了不同的单元格样式（红色、蓝色、绿色）。
设定了列宽、对齐方式、自动换行等属性，以优化视觉呈现和数据阅读性。
为整个工作表应用统一的风格设置，并设置冻结窗格以改善导航体验。
函数详解
NewXlsx：创建新的 Excel 文件并初始化工作表，包括定义列名和工作表的基础结构。
PutSht_INFO、PutSht_OS、PutSht_DB：
这些函数负责将相应的数据写入特定的工作表中。通过读取结构体中的数据，对应地填充到 Excel 的单元格中。
根据数据项的警报级别，应用不同的样式，如警告级别高的数据项会用红色标出，以便用户快速识别问题所在。
设计优化建议
错误处理：
在文件操作和数据转换过程中，代码应更全面地处理可能出现的错误。例如，当 excelize.OpenFile 或 strconv 函数调用失败时，应适当记录错误并处理，而不仅仅是打印错误信息。
性能考虑：
如果数据量大，每次写入单元格后都更新文件可能会影响性能。考虑在数据全部写入后再进行一次性保存。
代码重构：
多个 PutSht_ 函数中存在重复的代码逻辑，特别是在处理样式和写入单元格的部分。可以考虑抽象出共用的功能，如设置单元格内容和样式的函数，以减少代码冗余并提高可维护性。
这个模块是项目中处理数据输出的关键部分，良好的实现不仅可以提供清晰的数据展示，还能通过有效的错误处理和数据校验，提升整个应用的健壮性和用户体验


## chk.db
create table rules (id int PRIMARY KEY,nm varchar(10),desc varchar(30),type varchar(10),level int,rule text);
insert into rules values (1001,'osparameter','主机参数检查','OS',1,'查看主机内核参数配置情况;根据最佳实践建议: 
系统默认ASLR空间地址随机存在BUG隐患,根据ORACLE建议需要关闭ASLR ,修改randomze_va_space 1->0 或者2->0 （Doc ID 1345364.1）。
panic_on_oops存在BUG隐患，根据ORACLE建议设置该值为1。
kernel.panic_on_oops=1
vm.min_free_kbytes 内核参数最少保留 524288 (即512MB 注意单位K）以允许OS更更快地回收内存，可以避免内存低的压⼒。
net.ipv4.ip_local_port_range = 9000 65500，增加可用端口范围。
/etc/security/limits.conf 
* soft nproc 16384
* hard nproc 16384
* soft nofile 65536
* hard nofile 65536
注 : 达到SOFT限制应用也会报错,因此SOFT和HARD都需要检查
');



## Bug List
关于3.3.18序列最大值使用检查这一块的，当MAXVALUE为0时查询语句select sequence_owner,sequence_name, max_value,last_number,cache_size,round(last_number/max_value ,2) percent_use from dba_sequences 
where  last_number/max_value >0.8 and  cycle_flag='N'会报错ORA-01476: divisor is equal to zero，看看要不要加上max_value<>0