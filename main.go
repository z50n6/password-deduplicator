package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 显示欢迎横幅
	displayBanner()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("🔍 请输入密码字典文件路径 (例如：/path/to/passwd.txt): ")
		inputPath, _ := reader.ReadString('\n')
		inputPath = strings.TrimSpace(inputPath)

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			fmt.Println("❌ 文件不存在，请检查路径后重试。")
			continue
		}

		err := processFile(inputPath)
		if err != nil {
			fmt.Printf("⚠️ 处理文件过程中发生错误: %v\n", err)
			continue
		}

		fmt.Println("\033[32m✅ 处理完成!\033[0m") // 绿色输出

		fmt.Print("🔄 是否继续? (y继续, q退出): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "q" {
			fmt.Println("👋 退出程序。")
			break
		} else if choice != "y" {
			fmt.Println("❌ 无效选择，退出程序。")
			break
		}
	}
}

// 显示程序横幅
func displayBanner() {
	fmt.Println("===================================")
	fmt.Println("          字典去重工具        ")
	fmt.Println("              Version 1.0          ")
	fmt.Println("          欢迎使用本工具！         ")
	fmt.Println("===================================")
	fmt.Println()
}

func processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	lineCount := make(map[string]int)
	scanner := bufio.NewScanner(file)

	// 统计每一行出现的次数
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lineCount[line]++
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println("\033[34m🔄 重复的内容:\033[0m") // 蓝色输出
	hasDuplicates := false
	uniqueLines := []string{}

	// 输出重复内容并收集唯一行
	for line, count := range lineCount {
		if count > 1 {
			fmt.Printf("🔴 '%s' 出现了 %d 次\n", line, count)
			hasDuplicates = true
		}
		uniqueLines = append(uniqueLines, line)
	}

	if !hasDuplicates {
		fmt.Println("✅ 没有找到重复的内容。")
		return nil
	}

	// 将去重后的内容写回源文件
	return writeUniqueLines(filePath, uniqueLines)
}

func writeUniqueLines(filePath string, lines []string) error {
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	for _, line := range lines {
		_, err := outputFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	fmt.Println("\033[32m✅ 去重后的文件已更新。\033[0m") // 绿色输出
	return nil
}
