package privateservice

import (
	"fmt"
	gqaGlobal "github.com/Junvary/gin-quasar-admin/GQA-BACKEND/global"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/gqaplugin/example/model"
	gqaModel "github.com/Junvary/gin-quasar-admin/GQA-BACKEND/model"
	gqaServicePrivate "github.com/Junvary/gin-quasar-admin/GQA-BACKEND/service/private"
	gqaUtils "github.com/Junvary/gin-quasar-admin/GQA-BACKEND/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func GetTestDataList(getTestDataList model.RequestGetTestDataList, username string) (err error, exportData []model.GqaPluginExampleTestData, total int64) {
	pageSize := getTestDataList.PageSize
	offset := getTestDataList.PageSize * (getTestDataList.Page - 1)
	var exportDataList []model.GqaPluginExampleTestData
	var db *gorm.DB
	if err, db = gqaServicePrivate.DeptDataPermission(username, gqaGlobal.GqaDb.Model(&model.GqaPluginExampleTestData{})); err != nil {
		return err, exportDataList, 0
	}
	//配置搜索
	if getTestDataList.Column1 != "" {
		db = db.Where("column1 like ?", "%"+getTestDataList.Column1+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, exportDataList, 0
	}
	err = db.Limit(pageSize).Offset(offset).Order(gqaModel.OrderByColumn(getTestDataList.SortBy, getTestDataList.Desc)).Preload("CreatedByUser").Find(&exportDataList).Error
	return err, exportDataList, total
}

func EditTestData(toEditTestData model.GqaPluginExampleTestData, username string) (err error) {
	var db *gorm.DB
	if err, db = gqaServicePrivate.DeptDataPermission(username, gqaGlobal.GqaDb.Model(&model.GqaPluginExampleTestData{})); err != nil {
		return err
	}
	var exportData model.GqaPluginExampleTestData
	if err = db.Where("id = ?", toEditTestData.Id).First(&exportData).Error; err != nil {
		return err
	}
	//err = gqaGlobal.GqaDb.Updates(&toEditTestData).Error
	err = db.Updates(gqaUtils.MergeMap(gqaUtils.GlobalModelToMap(&toEditTestData.GqaModel),
		map[string]interface{}{
			"column1": toEditTestData.Column1,
			"column2": toEditTestData.Column2,
			"column3": toEditTestData.Column3,
			"column4": toEditTestData.Column4,
			"column5": toEditTestData.Column5,
		})).Error
	return err
}

func AddTestData(toAddTestData model.GqaPluginExampleTestData, username string) (err error) {
	var db *gorm.DB
	if err, db = gqaServicePrivate.DeptDataPermission(username, gqaGlobal.GqaDb.Model(&model.GqaPluginExampleTestData{})); err != nil {
		return err
	}
	err = db.Create(&toAddTestData).Error
	return err
}

func DeleteTestDataById(id uint, username string) (err error) {
	var db *gorm.DB
	if err, db = gqaServicePrivate.DeptDataPermission(username, gqaGlobal.GqaDb.Model(&model.GqaPluginExampleTestData{})); err != nil {
		return err
	}
	var exportData model.GqaPluginExampleTestData
	if err = db.Where("id = ?", id).First(&exportData).Error; err != nil {
		return err
	}
	err = db.Where("id = ?", id).Unscoped().Delete(&exportData).Error
	return err
}

func QueryTestDataById(id uint, username string) (err error, exportDataInfo model.GqaPluginExampleTestData) {
	var exportData model.GqaPluginExampleTestData
	var db *gorm.DB
	if err, db = gqaServicePrivate.DeptDataPermission(username, gqaGlobal.GqaDb.Model(&model.GqaPluginExampleTestData{})); err != nil {
		return err, exportData
	}
	err = db.Preload("CreatedByUser").Preload("UpdatedByUser").First(&exportData, "id = ?", id).Error
	return err, exportData
}

func ExportTestData(getTestDataList model.RequestGetTestDataList, filePath string, username string) error {
	excel := excelize.NewFile()
	err := excel.SetSheetRow("Sheet1", "A1", &[]string{"第1列", "第2列", "第3列", "第4列", "第5列"})
	if err != nil {
		return err
	}
	var db = gqaGlobal.GqaDb
	var exportDataList []model.GqaPluginExampleTestData
	//这里可以加入查询条件
	if getTestDataList.Column1 != "" {
		db = db.Where("column1 like ?", "%"+getTestDataList.Column1+"%")
	}
	err = db.Order(gqaModel.OrderByColumn(getTestDataList.SortBy, getTestDataList.Desc)).Find(&exportDataList).Error
	if err != nil {
		return err
	} else {
		for i, data := range exportDataList {
			key := fmt.Sprintf("A%d", i+2)
			err = excel.SetSheetRow("Sheet1", key, &[]interface{}{
				data.Column1,
				data.Column2,
				data.Column3,
				data.Column4,
				data.Column5,
			})
			if err != nil {
				return err
			}
		}
		err = excel.SaveAs(filePath)
		return err
	}
}
