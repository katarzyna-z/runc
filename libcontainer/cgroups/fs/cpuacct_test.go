// +build linux

package fs

import (
	"testing"

	"github.com/opencontainers/runc/libcontainer/cgroups"
)

const (
	cpuAcctUsageContents       = "12262454190222160"
	cpuAcctUsagePerCPUContents = "1564936537989058 1583937096487821 1604195415465681 1596445226820187 1481069084155629 1478735613864327 1477610593414743 1476362015778086"
	cpuAcctStatContents        = "user 452278264\nsystem 291429664"
)

func TestCpuacctStats(t *testing.T) {
	helper := NewCgroupTestUtil("cpuacct.", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"cpuacct.usage":        cpuAcctUsageContents,
		"cpuacct.usage_percpu": cpuAcctUsagePerCPUContents,
		"cpuacct.stat":         cpuAcctStatContents,
	})

	cpuacct := &CpuacctGroup{}
	actualStats := *cgroups.NewStats()
	err := cpuacct.GetStats(helper.CgroupPath, &actualStats)
	if err != nil {
		t.Fatal(err)
	}

	expectedStats := cgroups.CpuUsage{
		TotalUsage:        uint64(12262454190222160),
		PercpuUsage:       []uint64{1564936537989058, 1583937096487821, 1604195415465681, 1596445226820187, 1481069084155629, 1478735613864327, 1477610593414743, 1476362015778086},
		UsageInKernelmode: (uint64(291429664) * nanosecondsInSecond) / clockTicks,
		UsageInUsermode:   (uint64(452278264) * nanosecondsInSecond) / clockTicks,
	}

	expectCPUUsageStatsEqual(t, expectedStats, actualStats.CpuStats.CpuUsage)
}
