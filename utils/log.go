package utils

import "github.com/sealyun/lvscare/internal/klog"

func Config(level string) {
	switch level {
	case "INFO":
		klog.InitVerbosity(0)
	case "DEBG":
		klog.InitVerbosity(9)
	default:
		klog.InitVerbosity(0)
	}
}
