package avocados

import (
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"strconv"
)

//HTTPService ...
type HTTPService interface {
	Register (*gin.Engine)
}

type httpService struct {
	endpoints []*endpoint
}

type endpoint struct {
	method string
	path string
	function gin.HandlerFunc
}

//NewHTTPTransport ...
func NewHTTPTransport(s Service) HTTPService {
	endpoints:= makeEndpoints(s)
	return httpService{endpoints}
}

func makeEndpoints(s Service) []*endpoint{
	list := []*endpoint{}
	list = append(list, &endpoint{
		method : "GET",
		path : "/avocados",
		function: getAll(s),
	})
	list = append(list, &endpoint{
		method : "POST",
		path : "/avocados/add",
		function: addAvocado(s),
	})
	list = append(list, &endpoint{
		method : "GET",
		path : "/avocados/:id",
		function: getByID(s),
	})
	list = append(list, &endpoint{
		method : "DELETE",
		path : "/avocados/:id",
		function: deleteByID(s),
	})
	list = append(list, &endpoint{
		method : "POST",
		path : "/avocados/changequantity/:id",
		function: changeQuantity(s),
	})
	list = append(list, &endpoint{
		method : "POST",
		path : "/avocados/changeprice/:id",
		function: changePrice(s),
	})
	return list
}
func changeQuantity(s Service) gin.HandlerFunc{

	return func (c*gin.Context){
		var avoData Avocado
		body:= c.Request.Body
		x, _ := ioutil.ReadAll(body)
		ID, _ := strconv.ParseInt(c.Param("id"), 6, 12)
		json.Unmarshal([]byte(x), &avoData)
		
		if avoData.Quantity >= 0{
			res, err := s.ChangeAvoQuantity(ID, avoData.Quantity)
			c.JSON(http.StatusOK, gin.H{
				"response" : res,
				"error" : err,
			})
		}else{
			c.JSON(http.StatusConflict, gin.H{
				"error" : "Ingrese una cantidad...",
			})		
		}
	}
}

func changePrice(s Service) gin.HandlerFunc{

	return func (c*gin.Context){
		var avoData Avocado
		body:= c.Request.Body
		x, _ := ioutil.ReadAll(body)
		ID, _ := strconv.ParseInt(c.Param("id"), 6, 12)
		json.Unmarshal([]byte(x), &avoData)

		if avoData.Price >= 0{
			res, err := s.ChangeAvoPrice(ID, avoData.Price)
			c.JSON(http.StatusOK, gin.H{
				"response" : res,
				"error" : err,
			})
		}else{
			c.JSON(http.StatusConflict, gin.H{
				"error" : "Ingrese un valor...",
			})		
		}
	}
}

func getByID(s Service) gin.HandlerFunc{
	return func (c*gin.Context){
		ID, _ := strconv.ParseInt(c.Param("id"), 6, 12) 
		avo, err := s.GetByID(ID)
		c.JSON(http.StatusOK, gin.H{
			"avo": avo,
			"response" : err,
		})
	}
}

func deleteByID(s Service) gin.HandlerFunc{
	return func (c*gin.Context){
		ID, _ := strconv.ParseInt(c.Param("id"), 6, 12) 
		res, err := s.DeleteByID(ID)
		c.JSON(http.StatusOK, gin.H{
			"avo": res,
			"response": err,
		})
	}
}

func addAvocado(s Service) gin.HandlerFunc{
	return func(c *gin.Context){
		body := c.Request.Body
		x, _ := ioutil.ReadAll(body)
		var avoData Avocado
		json.Unmarshal([]byte(x), &avoData)
		avo := Avocado{0, avoData.Name, avoData.Image, avoData.Price, avoData.Stock, avoData.Quantity}

		if avo.Name != "" && avo.Image != "" && avo.Price >= 0{
			c.JSON(http.StatusCreated, gin.H{
				"avocados": s.AddAvocado(avo),
			})
		} else{
			c.JSON(http.StatusConflict, gin.H{
				"Error" : "Complete todos los datos...",
			})
		}
	}
}

func getAll(s Service) gin.HandlerFunc{
	return func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"avocado" : s.GetAll(),
			})
	}
}

func (s httpService) Register( r *gin.Engine){
	for _, e:= range s.endpoints {
		r.Handle(e.method, e.path, e.function)
	}
}