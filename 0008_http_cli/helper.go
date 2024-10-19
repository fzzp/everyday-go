package main

import "net/url"

// NOTE: 这里可以直接使用零依赖的govalidator库
// https://github.com/asaskevich/govalidator
// 下面的验证是有点瑕疵的
// 比如：https://xxx/https://...会被认为是正确的Url

// ValidateURL 验证字符串是否是URL地址
func ValidateURL(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	u, err := url.Parse(link)
	if err != nil {
		return false
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return false
	}

	return true
}

// findKeyByArgs 查询是否有 -c 命令，有则删除，并返回新的 args
func findKeyByArgs(args []string, key string) ([]string, bool) {
	var newArgs []string
	var cmdC = false
	for _, arg := range args {
		if arg == key {
			cmdC = true
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	return newArgs, cmdC
}
