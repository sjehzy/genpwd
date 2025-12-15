package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// 定义字符集
const (
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	digitChars   = "0123456789"
	symbolChars  = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	similarChars = "0OolI18B9q" // 易混淆字符
)

// 配置项
type Config struct {
	length         int  // 密码长度
	includeUpper   bool // 包含大写字母
	includeLower   bool // 包含小写字母
	includeDigits  bool // 包含数字
	includeSymbols bool // 包含特殊符号（默认关闭）
	excludeSimilar bool // 排除易混淆字符（默认开启，即不启用易混淆字符）
	batch          int  // 批量生成数量
	copyClipboard  bool // 复制到剪贴板（仅 Windows 支持）
}

// 初始化随机数生成器（修复：正确初始化随机数）
func init() {
	rand.Seed(time.Now().UnixNano()) // 兼容所有 Go 版本，避免 1.20+ 以下报错
}

// 过滤易混淆字符
func filterSimilarChars(chars string) string {
	var filtered strings.Builder
	for _, c := range chars {
		if !strings.ContainsRune(similarChars, c) {
			filtered.WriteRune(c)
		}
	}
	return filtered.String()
}

// 生成单个密码（逻辑不变，仅适配默认配置）
func generatePassword(cfg Config) (string, error) {
	// 构建字符池（严格按配置过滤）
	var charPool strings.Builder

	// 大写字母
	if cfg.includeUpper {
		chars := upperChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			charPool.WriteString(chars)
		}
	}
	// 小写字母
	if cfg.includeLower {
		chars := lowerChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			charPool.WriteString(chars)
		}
	}
	// 数字
	if cfg.includeDigits {
		chars := digitChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			charPool.WriteString(chars)
		}
	}
	// 特殊符号（默认关闭，需手动 -s true 启用）
	if cfg.includeSymbols {
		chars := symbolChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			charPool.WriteString(chars)
		}
	}

	// 校验字符池
	pool := charPool.String()
	if len(pool) == 0 {
		return "", fmt.Errorf("字符池为空，请至少启用一种字符类型（-u/-w/-d/-s）")
	}
	if cfg.length < 1 {
		return "", fmt.Errorf("密码长度必须 ≥ 1（建议 ≥ 8）")
	}

	// 确保每种选中类型至少包含 1 个字符
	var password []rune

	// 大写字母
	if cfg.includeUpper {
		chars := upperChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			password = append(password, rune(chars[rand.Intn(len(chars))]))
		}
	}
	// 小写字母
	if cfg.includeLower {
		chars := lowerChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			password = append(password, rune(chars[rand.Intn(len(chars))]))
		}
	}
	// 数字
	if cfg.includeDigits {
		chars := digitChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			password = append(password, rune(chars[rand.Intn(len(chars))]))
		}
	}
	// 特殊符号（仅启用时添加）
	if cfg.includeSymbols {
		chars := symbolChars
		if cfg.excludeSimilar {
			chars = filterSimilarChars(chars)
		}
		if len(chars) > 0 {
			password = append(password, rune(chars[rand.Intn(len(chars))]))
		}
	}

	// 补充剩余长度
	remaining := cfg.length - len(password)
	if remaining > 0 {
		poolRunes := []rune(pool)
		for i := 0; i < remaining; i++ {
			password = append(password, poolRunes[rand.Intn(len(poolRunes))])
		}
	}

	// 打乱密码顺序
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password), nil
}

// 复制到剪贴板（Windows 稳定版）
func copyToClipboard(text string) error {
	escapedText := strings.ReplaceAll(text, "'", "''")
	cmd := exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("Set-Clipboard -Value '%s'", escapedText))
	return cmd.Run()
}

func main() {
	// 核心调整：默认配置改为「仅大小写+数字，关闭特殊符号+易混淆字符」
	var cfg Config
	// 密码长度默认12
	flag.IntVar(&cfg.length, "l", 12, "密码长度（默认12），等价 -len")
	flag.IntVar(&cfg.length, "len", 12, "【兼容】密码长度（默认12）")

	// 大写字母：默认启用
	flag.BoolVar(&cfg.includeUpper, "u", true, "包含大写字母（默认true），等价 -upper")
	flag.BoolVar(&cfg.includeUpper, "upper", true, "【兼容】包含大写字母（默认true）")

	// 小写字母：默认启用
	flag.BoolVar(&cfg.includeLower, "w", true, "包含小写字母（默认true），等价 -lower（w=lower）")
	flag.BoolVar(&cfg.includeLower, "lower", true, "【兼容】包含小写字母（默认true）")

	// 数字：默认启用
	flag.BoolVar(&cfg.includeDigits, "d", true, "包含数字（默认true），等价 -digit")
	flag.BoolVar(&cfg.includeDigits, "digit", true, "【兼容】包含数字（默认true）")

	// 特殊符号：默认关闭（需手动 -s true 启用）
	flag.BoolVar(&cfg.includeSymbols, "s", false, "包含特殊符号（默认false），等价 -symbol")
	flag.BoolVar(&cfg.includeSymbols, "symbol", false, "【兼容】包含特殊符号（默认false）")

	// 排除易混淆字符：默认开启（需手动 -n false 禁用，即启用易混淆字符）
	flag.BoolVar(&cfg.excludeSimilar, "n", true, "排除易混淆字符（默认true），等价 -no-similar；-n false 启用易混淆字符")
	flag.BoolVar(&cfg.excludeSimilar, "no-similar", true, "【兼容】排除易混淆字符（默认true）")

	// 批量生成：默认1个
	flag.IntVar(&cfg.batch, "b", 1, "批量生成数量（默认1），等价 -batch")
	flag.IntVar(&cfg.batch, "batch", 1, "【兼容】批量生成数量（默认1）")

	// 复制剪贴板：默认关闭
	flag.BoolVar(&cfg.copyClipboard, "c", false, "复制第一个密码到剪贴板（默认false），等价 -copy")
	flag.BoolVar(&cfg.copyClipboard, "copy", false, "【兼容】复制第一个密码到剪贴板（默认false）")

	// 自定义帮助信息（明确默认规则）
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使用方法：%s [参数]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "默认规则：仅生成 大小写字母+数字（无特殊符号、无易混淆字符）")
		fmt.Fprintln(os.Stderr, "参数说明（推荐首字母缩写）：")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\n示例：")
		fmt.Fprintln(os.Stderr, "  1. 默认生成12位密码（大小写+数字，无特殊符号/易混淆字符）：", os.Args[0])
		fmt.Fprintln(os.Stderr, "  2. 生成16位包含特殊符号的密码：", os.Args[0], "-l 16 -s true")
		fmt.Fprintln(os.Stderr, "  3. 生成8位包含易混淆字符的纯数字：", os.Args[0], "-l 8 -u false -w false -n false")
		fmt.Fprintln(os.Stderr, "  4. 生成20位包含特殊符号+易混淆字符的密码：", os.Args[0], "-l 20 -s true -n false")
	}
	flag.Parse()

	// 批量生成密码
	var passwords []string
	for i := 0; i < cfg.batch; i++ {
		pwd, err := generatePassword(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ 生成密码失败：%v\n", err)
			os.Exit(1)
		}
		passwords = append(passwords, pwd)
	}

	// 输出结果
	fmt.Println("=== 生成的随机密码 ===")
	for i, pwd := range passwords {
		fmt.Printf("%d. %s\n", i+1, pwd)
	}

	// 复制到剪贴板
	if cfg.copyClipboard && len(passwords) > 0 {
		if err := copyToClipboard(passwords[0]); err != nil {
			fmt.Fprintf(os.Stderr, "❌ 复制到剪贴板失败：%v\n", err)
		} else {
			fmt.Println("\n✅ 第一个密码已复制到剪贴板！")
		}
	}
}
