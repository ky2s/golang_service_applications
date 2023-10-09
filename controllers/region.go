package controllers

import (
	"fmt"
	"net/http"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// interface
type RegionController interface {
	GetProvince(c *gin.Context)
	GetCity(c *gin.Context)
	GetDistrict(c *gin.Context)
	GetSubDistrict(c *gin.Context)
	GetRadius(c *gin.Context)
}

type regionController struct {
	regionMod models.RegionModels
}

func NewRegionController(regionModel models.RegionModels) RegionController {
	return &regionController{
		regionMod: regionModel,
	}
}

func (ctr *regionController) GetProvince(c *gin.Context) {

	getProvinces, err := ctr.regionMod.GetListProvinces(tables.Province{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	if len(getProvinces) > 0 {
		var result []objects.Province
		for i := 0; i < len(getProvinces); i++ {

			var each objects.Province
			each.ID = getProvinces[i].ID
			each.Name = getProvinces[i].Name
			each.Status = getProvinces[i].Status
			each.CreatedAt = getProvinces[i].CreatedAt
			each.UpdatedAt = getProvinces[i].UpdatedAt
			each.DeletedAt = getProvinces[i].DeletedAt

			result = append(result, each)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data is available",
			"data":    result,
		})
		return
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"status":  false,
			"message": "Data is not available",
			"data":    nil,
		})
		return
	}
}

func (ctr *regionController) GetCity(c *gin.Context) {

	provinceID, _ := strconv.Atoi(c.Param("provinceid"))

	if provinceID > 0 {
		//jika ada province id
		var where tables.Cities
		where.ProvinceID = provinceID

		getCities, err := ctr.regionMod.GetListCitiesByProvinceID(where)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(getCities) > 0 {
			var result []objects.Cities
			for i := 0; i < len(getCities); i++ {

				var each objects.Cities
				each.ID = getCities[i].ID
				each.ProvinceID = getCities[i].ProvinceID
				each.Name = getCities[i].Name
				each.Status = getCities[i].Status
				each.CreatedAt = getCities[i].CreatedAt
				each.UpdatedAt = getCities[i].UpdatedAt
				each.DeletedAt = getCities[i].DeletedAt

				result = append(result, each)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    result,
			})
			return
		} else {
			c.JSON(http.StatusNoContent, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}

	} else {
		//jika tidak ada province id
		getCities, err := ctr.regionMod.GetListAllCities(tables.Cities{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(getCities) > 0 {
			var result []objects.Cities
			for i := 0; i < len(getCities); i++ {

				var each objects.Cities
				each.ID = getCities[i].ID
				each.ProvinceID = getCities[i].ProvinceID
				each.Name = getCities[i].Name
				each.Status = getCities[i].Status
				each.CreatedAt = getCities[i].CreatedAt
				each.UpdatedAt = getCities[i].UpdatedAt
				each.DeletedAt = getCities[i].DeletedAt

				result = append(result, each)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    result,
			})
			return
		} else {
			c.JSON(http.StatusNoContent, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	}
}

func (ctr *regionController) GetDistrict(c *gin.Context) {

	cityID, _ := strconv.Atoi(c.Param("cityid"))

	if cityID > 0 {
		//jika ada province id
		var where tables.Districts
		where.CityID = cityID

		getDistrict, err := ctr.regionMod.GetListDistrictsByProvinceID(where)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(getDistrict) > 0 {
			var result []objects.Districts
			for i := 0; i < len(getDistrict); i++ {

				var each objects.Districts
				each.ID = getDistrict[i].ID
				each.CityID = getDistrict[i].CityID
				each.Name = getDistrict[i].Name
				each.Status = getDistrict[i].Status
				each.CreatedAt = getDistrict[i].CreatedAt
				each.UpdatedAt = getDistrict[i].UpdatedAt
				each.DeletedAt = getDistrict[i].DeletedAt

				result = append(result, each)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    result,
			})
			return
		} else {
			c.JSON(http.StatusNoContent, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}

	} else {
		//jika tidak ada province id
		getDistrict, err := ctr.regionMod.GetListAllDistricts(tables.Districts{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(getDistrict) > 0 {
			var result []objects.Districts
			for i := 0; i < len(getDistrict); i++ {

				var each objects.Districts
				each.ID = getDistrict[i].ID
				each.CityID = getDistrict[i].CityID
				each.Name = getDistrict[i].Name
				each.Status = getDistrict[i].Status
				each.CreatedAt = getDistrict[i].CreatedAt
				each.UpdatedAt = getDistrict[i].UpdatedAt
				each.DeletedAt = getDistrict[i].DeletedAt

				result = append(result, each)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    result,
			})
			return
		} else {
			c.JSON(http.StatusNoContent, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	}
}

func (ctr *regionController) GetSubDistrict(c *gin.Context) {

	districtID, _ := strconv.Atoi(c.Param("districtid"))

	if districtID > 0 {
		//jika ada province id
		var where tables.SubDistricts
		where.DistrictID = districtID

		getSubDistrict, err := ctr.regionMod.GetListSubDistrictsByProvinceID(where)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(getSubDistrict) > 0 {
			var result []objects.SubDistricts
			for i := 0; i < len(getSubDistrict); i++ {

				var each objects.SubDistricts
				each.ID = getSubDistrict[i].ID
				each.DistrictID = getSubDistrict[i].DistrictID
				each.Name = getSubDistrict[i].Name
				each.PostalCode = getSubDistrict[i].PostalCode
				each.Status = getSubDistrict[i].Status
				each.CreatedAt = getSubDistrict[i].CreatedAt
				each.UpdatedAt = getSubDistrict[i].UpdatedAt
				each.DeletedAt = getSubDistrict[i].DeletedAt

				result = append(result, each)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    result,
			})
			return
		} else {
			c.JSON(http.StatusNoContent, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}

	} else {
		//jika tidak ada province id
		getSubDistrict, err := ctr.regionMod.GetListAllSubDistricts(tables.SubDistricts{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		if len(getSubDistrict) > 0 {
			var result []objects.SubDistricts
			for i := 0; i < len(getSubDistrict); i++ {

				var each objects.SubDistricts
				each.ID = getSubDistrict[i].ID
				each.DistrictID = getSubDistrict[i].DistrictID
				each.Name = getSubDistrict[i].Name
				each.PostalCode = getSubDistrict[i].PostalCode
				each.Status = getSubDistrict[i].Status
				each.CreatedAt = getSubDistrict[i].CreatedAt
				each.UpdatedAt = getSubDistrict[i].UpdatedAt
				each.DeletedAt = getSubDistrict[i].DeletedAt

				result = append(result, each)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    result,
			})
			return
		} else {
			c.JSON(http.StatusNoContent, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	}
}

func (ctr *regionController) GetRadius(c *gin.Context) {

	var reqData objects.Radius
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	var radius objects.Radius
	radius.LocationID = reqData.LocationID
	radius.Latitude = reqData.Latitude
	radius.Longitude = reqData.Longitude

	get, err := ctr.regionMod.GetRadius(radius)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	var radiusshow objects.RadiusSlow
	radiusshow.Distance = get.Distance
	radiusshow.IsRadius = get.IsRadius

	if get.IsRadius == true {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "User sudah dalam radius",
			"data":    radiusshow,
		})
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User di luar radius",
			"data":    radiusshow,
		})
		return

	}

}
