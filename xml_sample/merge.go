package main

import (
	"fmt"
	"log"

	"github.com/beevik/etree"
)

func mergeXML(file1, file2, outputFile string) {
	doc1 := etree.NewDocument()
	if err := doc1.ReadFromFile(file1); err != nil {
		log.Fatal(err)
	}

	doc2 := etree.NewDocument()
	if err := doc2.ReadFromFile(file2); err != nil {
		log.Fatal(err)
	}

	// 创建输出文档和根元素
	outputDoc := etree.NewDocument()
	outputDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := outputDoc.CreateElement("EACHK")

	// 处理 TAG0
	tag0 := root.CreateElement("TAG0")
	node1 := tag0.CreateElement("NODE1")
	for _, elem := range doc1.FindElement("./EACHK/TAG0").ChildElements() {
		node1.AddChild(elem.Copy())
	}
	node2 := tag0.CreateElement("NODE2")
	for _, elem := range doc2.FindElement("./EACHK/TAG0").ChildElements() {
		node2.AddChild(elem.Copy())
	}

	// 处理 TAG1 (仅使用第一个文件的内容)
	tag1 := root.CreateElement("TAG1")
	for _, elem := range doc1.FindElement("./EACHK/TAG1").ChildElements() {
		tag1.AddChild(elem.Copy())
	}

	// 处理 TAG2
	tag2 := root.CreateElement("TAG2")
	node1Tag2 := tag2.CreateElement("NODE1")
	for _, elem := range doc1.FindElement("./EACHK/TAG2").ChildElements() {
		node1Tag2.AddChild(elem.Copy())
	}
	node2Tag2 := tag2.CreateElement("NODE2")
	for _, elem := range doc2.FindElement("./EACHK/TAG2").ChildElements() {
		node2Tag2.AddChild(elem.Copy())
	}

	// 保存合并后的XML到文件
	outputDoc.Indent(2)
	if err := outputDoc.WriteToFile(outputFile); err != nil {
		log.Fatal(err)
	}

	fmt.Println("XML文件合并完成：", outputFile)
}

func main() {
	file1 := "20240401_M001R01DG_lm0ora01_M001R011.xml"
	file2 := "20240401_M001R01DG_lm0ora02_M001R012.xml"
	outputFile := "Merged_RAC.xml"

	mergeXML(file1, file2, outputFile)
}
