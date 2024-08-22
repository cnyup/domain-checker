package checker

import (
	"fmt"
)

// Check2StdoutByFile
//
//	@Description: 单文件检测结果控制台打印
//	@receiver c
//	@param file
//	@param days
//	@return error
func (c *Checker) Check2StdoutByFile(file string, days int) error {
	err := c.CheckByFile(file, days, c.stdoutInfo)
	if err != nil {
		return err
	}
	return nil
}

// Check2StdoutByDir
//
//	@Description:  按照后缀检索目录文件获取检测结果并控制台打印
//	@receiver c
//	@param dir
//	@param su
//	@param days
//	@return error
func (c *Checker) Check2StdoutByDir(dir, su string, days int) error {
	err := c.CheckByDir(dir, su, days, c.stdoutInfo)
	if err != nil {
		return err
	}
	return nil
}

// stdoutInfo
//
//	@Description: 控制台打印
//	@receiver c
func (c *Checker) stdoutInfo() {
	fmt.Printf("hostname is %s\n", c.EcsInfo.Name)
	fmt.Printf("wanip is %s\n", c.EcsInfo.WanIp)
	fmt.Printf("lanip is %s\n", c.EcsInfo.LanIp)

	if len(c.ExpireDomain) != 0 {
		fmt.Printf("expired domain: ")
		for _, info := range c.ExpireDomain {
			fmt.Printf("%s\t", info.DomainName)
		}
		fmt.Println()
	}
	if len(c.ThresholdDomain) != 0 {
		fmt.Println("Threshold domain: ")
		for _, info := range c.ThresholdDomain {
			fmt.Printf("%s have %d days later will expire\n", info.DomainName, info.ExpiredDays)
		}
	}

}
