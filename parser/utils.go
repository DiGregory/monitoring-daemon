package parser

import (
	"strings"
	"os/exec"
	"strconv"
	"regexp"
	"math"
	"errors"
)

var nilOutputErr = errors.New("server returned nothing")

var floatNumRegExp = regexp.MustCompile(`\d+,\d+`)
var integerNumRegExp = regexp.MustCompile(`\s\d+`)
var textRegExp = regexp.MustCompile(`\S+`)

//разбить команду на имя и аргументы
func parseRequestParams(request string) (command string, args []string) {
	r := strings.Split(request, " ")
	return r[0], r[1:]
}

//выполнить команду
func execCommand(command string, args []string) (*string, error) {
	dataCmd, err := exec.Command(command, args...).Output()
	if err != nil {
		return nil, err
	}
	if dataCmd == nil {
		return nil, nilOutputErr
	}
	output := string(dataCmd)
	return &output, nil
}

//захват чисел из строки значений.
func parseDataRow(rawValues string, re *regexp.Regexp) (loadValues []float32, err error) {
	rawString := re.FindAllString(rawValues, -1)
	for _, l := range rawString {
		stringWithDot := strings.Replace(l, ",", ".", -1)
		stringWithDot = strings.TrimSpace(stringWithDot)
		loadValue, err := strconv.ParseFloat(stringWithDot, 32)
		if err != nil {
			return nil, err
		}
		loadValues = append(loadValues, float32(loadValue))
	}
	return loadValues, nil
}

//имена строк
func parseDataFrameNames(dataFrame string, re *regexp.Regexp, rowIndex int, offset int, columnIndices []int) ([][]string, error) {
	rawRows := strings.Split(dataFrame, "\n")
	if offset == -1 {
		offset = len(rawRows) - rowIndex
	}
	var sortedDataFrameNames [][]string
	for i, v := range columnIndices {
		sortedDataFrameNames = append(sortedDataFrameNames, []string{})
		for j := rowIndex; j < rowIndex+offset; j++ {
			names := re.FindAllString(rawRows[j], -1)
			if v < (len(names)) {
				sortedDataFrameNames[i] = append(sortedDataFrameNames[i], names[v])
			}
		}
	}

	return sortedDataFrameNames, nil
}

//захват значений из таблицы значений. rowIndex-с какой строки заканчивается заголовок, rowIndex+offset-на какой строке
//заканчиваются данные, columnIndices-нужные столбцы
func parseDataFrame(dataFrame string, re *regexp.Regexp, rowIndex int, offset int, columnIndices []int) ([][]float32, error) {
	rawRows := strings.Split(dataFrame, "\n")
	if offset == -1 {
		offset = len(rawRows) - rowIndex
	}
	var sortedDataFrame [][]float32
	for i, v := range columnIndices {
		sortedDataFrame = append(sortedDataFrame, []float32{})
		for j := rowIndex; j < rowIndex+offset; j++ {
			currentRowValues, err := parseDataRow(rawRows[j], re)
			if err != nil {
				return nil, err
			}
			sortedDataFrame[i] = append(sortedDataFrame[i], currentRowValues[v])
		}
	}
	return sortedDataFrame, nil
}

func round(x float32) float32 {

	return float32(math.Floor(float64(x*100)) / 100)
}

//сложения и деления структур
func sumLoadAverage(la1 *LoadAverage, la2 *LoadAverage) (la LoadAverage) {
	if la1 == nil || la2 == nil {
		return
	}
	return LoadAverage{la1.Load1 + la2.Load1,
		la1.Load5 + la2.Load5,
		la1.Load15 + la2.Load15}

}
func divisionLoadAverage(a *LoadAverage, m int) (la LoadAverage) {
	if a == nil {
		return
	}
	return LoadAverage{round(a.Load1 / float32(m)),
		round(a.Load5 / float32(m)),
		round(a.Load15 / float32(m))}

}

func sumCpuLoad(cl1 *CpuLoad, cl2 *CpuLoad) (cl CpuLoad) {
	if cl1 == nil || cl2 == nil {
		return
	}
	return CpuLoad{cl1.UserLoad + cl2.UserLoad,
		cl1.SystemLoad + cl2.SystemLoad,
		cl1.IDLELoad + cl2.IDLELoad}

}
func divisionCpuLoad(cl1 *CpuLoad, m int) (cl CpuLoad) {
	if cl1 == nil {
		return
	}
	return CpuLoad{round(cl1.UserLoad / float32(m)),
		round(cl1.SystemLoad / float32(m)),
		round(cl1.IDLELoad / float32(m))}

}

func sumDiskLoad(load1 *DiskLoad, load2 *DiskLoad) (sLoad DiskLoad) {
	if load1 == nil || load2 == nil {
		return
	}
	//если слайс пуст, заполнить нулями
	if len(load1.TPS) == 0 {
		for range load2.TPS {
			load1.TPS = append(load1.TPS, 0)
			load1.ReadLoad = append(load1.ReadLoad, 0)
			load1.WriteLoad = append(load1.WriteLoad, 0)
		}
	}
	sLoad.Names = append(sLoad.Names, load2.Names...)
	cpu := sumCpuLoad(&load1.CPU, &load2.CPU)
	sLoad.CPU = cpu
	for i := range load2.TPS {
		sLoad.TPS = append(sLoad.TPS, load1.TPS[i]+load2.TPS[i])
		sLoad.ReadLoad = append(sLoad.ReadLoad, load1.ReadLoad[i]+load2.ReadLoad[i])
		sLoad.WriteLoad = append(sLoad.WriteLoad, load1.WriteLoad[i]+load2.WriteLoad[i])
	}
	return sLoad
}

func divisionDiskLoad(load *DiskLoad, m int) (dLoad DiskLoad) {
	if load == nil {
		return
	}
	cpu := divisionCpuLoad(&load.CPU, m)
	load.CPU = cpu
	for i := range load.TPS {
		load.TPS[i] = round(load.TPS[i] / float32(m))
		load.WriteLoad[i] = round(load.WriteLoad[i] / float32(m))
		load.ReadLoad[i] = round(load.ReadLoad[i] / float32(m))
	}
	return *load
}

func sumDiskFree(free1 *DiskFree, free2 *DiskFree) (sFree DiskFree) {
	if free1 == nil || free2 == nil {
		return
	}
	//если слайс пуст, заполнить нулями
	if len(free1.INodeUses) == 0 {
		for range free2.INodeUses {
			free1.INodeUses = append(free1.INodeUses, 0)
			free1.INodeFree = append(free1.INodeFree, 0)
			free1.MBUses = append(free1.MBUses, 0)
			free1.MBFree = append(free1.MBFree, 0)
		}
	}
	sFree.Names = append(sFree.Names, free2.Names...)
	for i := range free2.MBUses {
		sFree.INodeUses = append(sFree.INodeUses, free1.INodeUses[i]+free2.INodeUses[i])
		sFree.MBUses = append(sFree.MBUses, free1.MBUses[i]+free2.MBUses[i])
		sFree.INodeFree = append(sFree.INodeFree, free1.INodeFree[i]+free2.INodeFree[i])
		sFree.MBFree = append(sFree.MBFree, free1.MBFree[i]+free2.MBFree[i])
	}
	return sFree
}

func divisionDiskFree(free *DiskFree, m int) (dFree DiskFree) {
	if free == nil {
		return
	}
	for i := range free.MBUses {
		free.MBUses[i] = round(free.MBUses[i] / float32(m))
		free.INodeUses[i] = round(free.INodeUses[i] / float32(m))
		free.MBFree[i] = round(free.MBFree[i] / float32(m))
		free.INodeFree[i] = round(free.INodeFree[i] / float32(m))
	}
	return *free
}

//получение усредненной статистики за m секунд
func GetAverageInfo(pd []ParsedData, m int) ParsedData {
	var (
		la      LoadAverage
		cl      CpuLoad
		df      DiskFree
		dl      DiskLoad
		average ParsedData
	)
	for _, v := range pd[len(pd)-m:] {
		la = sumLoadAverage(&la, &v.LoadAverage)
		cl = sumCpuLoad(&cl, &v.CpuLoad)
		df = sumDiskFree(&df, &v.DiskFree)
		dl = sumDiskLoad(&dl, &v.DiskLoad)
	}
	average.LoadAverage = divisionLoadAverage(&la, m)
	average.CpuLoad = divisionCpuLoad(&cl, m)
	average.DiskLoad = divisionDiskLoad(&dl, m)
	average.DiskFree = divisionDiskFree(&df, m)
	average.TopTalkers = pd[len(pd)-1].TopTalkers
	average.NetStat = pd[len(pd)-1].NetStat

	return average
}
