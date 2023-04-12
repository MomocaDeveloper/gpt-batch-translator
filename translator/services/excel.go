package services

import (
	"fmt"
	"strings"
	"github.com/xuri/excelize/v2"
)

type Cell struct{
	Id int 
	Source string 
	Target string 
	Require string
	Role string
	SecTarget string
}

type Content struct{
	Id int
	Require string
	Role string
	Remind string
	Cells []Cell
}

type ExcelObj struct{
	Name string
	SourceCol string
	TargetCol string
	RequireCol string
	RoleCol string
	RoleMap map[string]int
	SecTargetCol string
	Contents []*Content
}

type characterDes struct{
    roleDes string
    roleTranslateDes string
}

func ParseFullFile(file string) (*ExcelObj, error){
	filePath := fmt.Sprintf("uploads/%s", file)
	f, err := excelize.OpenFile(filePath)
	if err != nil{
		fmt.Println("excel open file fail", err)
		return &ExcelObj{}, err
	}
	resp := &ExcelObj{
		Name:file,
		RoleMap : make(map[string]int),
	}

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, rowErr := f.GetRows(sheet)
		defer f.Close()
		if rowErr != nil{
			fmt.Println("excel read from file rows fail", rowErr)
			return &ExcelObj{}, rowErr
		}

		for coln:=1; coln<=24; coln++{
			colName, _ := excelize.ColumnNumberToName(coln)
			cellVal, _ := f.GetCellValue(sheet, fmt.Sprintf("%s1", colName))
			if cellVal == ""{
				continue
			}
			if strings.Contains(cellVal, "原文"){ 
				if resp.SourceCol == ""{
					resp.SourceCol = colName
				}   
			}   
			if strings.Contains(cellVal, "输出"){ 
				if resp.TargetCol == ""{
					resp.TargetCol = colName
				}   
			}   
			if strings.Contains(cellVal, "要求"){ 
				if resp.RequireCol == ""{
					resp.RequireCol = colName
				}   
			}
			if strings.Contains(cellVal, "角色"){
				if resp.RoleCol == ""{
					resp.RoleCol = colName
				}
			}
			if strings.Contains(cellVal, "二次修改"){
				if resp.SecTargetCol == ""{
					resp.SecTargetCol = colName
				}
			}
		}

		
		
		for i, row := range rows{
			if len(row) == 0{
				continue
			}
			fmt.Println(len(row), row, i)
			if i==0{
				continue
			}
			source,_ := f.GetCellValue("sheet1", fmt.Sprintf("%s%d", resp.SourceCol, i+1))
			if len(source)==0{
				continue
			}
			target,_ := f.GetCellValue("sheet1", fmt.Sprintf("%s%d", resp.TargetCol, i+1))
			if len(target)!=0{
				continue
			}
			require,_ := f.GetCellValue("sheet1", fmt.Sprintf("%s%d", resp.RequireCol, i+1))
			role,_ := f.GetCellValue("sheet1", fmt.Sprintf("%s%d", resp.RoleCol, i+1))
			secTarget,_ := f.GetCellValue("sheet1", fmt.Sprintf("%s%d", resp.SecTargetCol, i+1))

			curLength := len(resp.Contents)
			curInd, exist := resp.RoleMap[role]
			if !exist||len(resp.Contents[curInd].Cells)>=20 {
				resp.RoleMap[role] = curLength
				resp.Contents = append(resp.Contents, &Content{
					Id : i+1,
					Role : role,
					Require : require,
				})
				curInd = curLength
				curLength += 1	
			}
			tmpCell := Cell{
				i+1,source,target,require,role,secTarget,
			}
			resp.Contents[curInd].Cells = append(resp.Contents[curInd].Cells, tmpCell)
		}
	}
	return resp, nil
}

func WriteIntoFile(obj ExcelObj) error {
	f := excelize.NewFile()
	outputFile := fmt.Sprintf("downloads/tr_%s", obj.Name)
	defer f.Close()
	sheetId := "sheet1"
	
	f.SetCellValue(sheetId, fmt.Sprintf("%s1", obj.SourceCol), "原文")
	f.SetCellValue(sheetId, fmt.Sprintf("%s1", obj.TargetCol), "输出")
	f.SetCellValue(sheetId, fmt.Sprintf("%s1", obj.RequireCol), "翻译要求")
	f.SetCellValue(sheetId, fmt.Sprintf("%s1", obj.RoleCol), "角色")
	f.SetCellValue(sheetId, fmt.Sprintf("%s1", obj.SecTargetCol), "二次修改")

	for _, content := range obj.Contents{
		for _, cell := range content.Cells{
			id := cell.Id

			source := cell.Source
			sourceId := fmt.Sprintf("%s%d", obj.SourceCol, id)
			f.SetCellValue(sheetId, sourceId, source)

			target := cell.Target
			targetId := fmt.Sprintf("%s%d", obj.TargetCol, id)
			f.SetCellValue(sheetId, targetId, target)

			require := cell.Require
			requireId := fmt.Sprintf("%s%d", obj.RequireCol, id)
			f.SetCellValue(sheetId, requireId, require)
			
			role := cell.Role
			roleId := fmt.Sprintf("%s%d", obj.RoleCol, id)
			f.SetCellValue(sheetId, roleId, role)

			secTarget := cell.SecTarget
			secTargetId := fmt.Sprintf("%s%d", obj.SecTargetCol, id)
			f.SetCellValue(sheetId, secTargetId, secTarget)
			fmt.Println("WriteIntoFile:", id, secTargetId, secTarget)
		}
	}
	if err := f.SaveAs(outputFile); err != nil {
		fmt.Println("write info file fail", err)
		return err
	}	   

	fmt.Println("save into output file success~")	   
	return nil
}

func (e ExcelObj)CalcTotalLine()int{
	total := 0
	for _, content := range e.Contents{
		total += len(content.Cells)
	}
	return total
}

func GetBaseNoun(file string)(map[string]string, error){
	filePath := fmt.Sprintf("keywords/%s", file)
	f, openErr := excelize.OpenFile(filePath)
	resultMap := make(map[string]string)
	if openErr != nil{
		fmt.Println("read noun from excel fail", file, openErr)
		return nil, openErr
	}
	sheets := f.GetSheetMap()
	for _, sheet := range sheets{
		rows, rowErr := f.GetRows(sheet)
		if rowErr != nil{
			fmt.Println("read file rows fail", rowErr)
			return resultMap, rowErr
		}

		for _, row := range rows{
			if len(row) == 0{
				continue
			}
			resultMap[row[0]] = row[1]
		}
	}
	return resultMap, nil
}

func GetBaseCharacter(file string)(map[string]characterDes, error){
	filePath := fmt.Sprintf("characters/%s", file)
	f, openErr := excelize.OpenFile(filePath)
	resultMap := make(map[string]characterDes)
	if openErr != nil{
		fmt.Println("read character from excel fail", file, openErr)
		return resultMap, openErr
	}   
	sheets := f.GetSheetMap()
	for _, sheet := range sheets{
		rows, rowErr := f.GetRows(sheet)
		if rowErr != nil{
			fmt.Println("read character file rows fail", rowErr)
			return resultMap, rowErr
		}   

		for _, row := range rows{
			if len(row) == 0{
				continue
			}   
			if len(row[1])==0{
				continue
			}
			transDesc := ""
			if len(row) >=3 && len(row[2])!=0{
				transDesc = row[2]
				fmt.Println("read character check:", row[0], transDesc)
			}
			c := characterDes{
                roleDes: row[1],
                roleTranslateDes: transDesc,
            }
            resultMap[row[0]] = c
		}   
	}   
	return resultMap, nil
}
