package user

import (
	"foodmap/internal/entity/user"
	"foodmap/internal/infra/delivery"
	"foodmap/internal/infra/object"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupAPI add api routes to server
func SetupAPI(s *delivery.Server, u user.IUserUsecase) {
	users := s.API.Group("/users")
	users.GET("/", getAPIUsers(u))
	users.POST("/", postAPIUsers(u))
	users.GET("/:id", getAPIUser(u))
	users.PUT("/:id", putAPIUser(u))
	users.DELETE("/:id", deleteAPIUser(u))
}

// getAPIUsers find a list of Users
//
// GET /api/Users
//
// query:
//   query (string) text to search with
//   limit (int64) limit the number of returned records, ignored when limit and
//     skip are both 0
//   skip (int64) records to skip from search results
//   fields ([]string) fields to return, id is always returned
//
// success response:
//   {
//      "data": {
//         "Users": [
//            {} // a user document, refer to getAPIUser
//         ]
//    	}
//   }
//
func getAPIUsers(u user.IUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, s, err := delivery.ParseLimitAndSkip(c.Query("limit"), c.Query("skip"))
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
		}
		result, err := u.Find(c.Query("query"), c.Query("fields"), l, s)
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"Users": result,
			},
		})
	}
}

// postAPIUsers create a store record
//
// POST /api/Users
//
// body: (json)
//   {
//      "name": "", // required, store display name, max length is 50
//      "description": "", // optional, store description, max length is 1000
//      "business_hours": [ // required, store business hours
//         {
//            "day": [ from_day, to_day ] // [2]int, required, accept monday (1) ~ sunday (7)
//            "time": [ from_time, to_time ] // [2]string, required, accept 00:00 ~ 23:59
//         }
//      ],
//      "categories": [], // []string, optional, required, categories of the store
//      "price_level": ""// required, must be one of the values: "cheap", "medium", "expensive"
//      "menu": [
//         {
//            "name": "", // required, product display name, max length is 50
//            "description": "", // optional, product description, max length is 1000
//            "category": "" // required, product category, max length is 50,
//            "price": 0, // int, *ignored/required* when the product *has/has no* variant,
//            "variants": [
//               {
//                  "name": "", // required, variant display name, max length is 5
//                  "price": 0 // int, required, variant price
//               }
//            ]
//         }
//      ]
//   }
//
// success response:
//   {
//      "id": "" // store id
//   }
//
func postAPIUsers(u user.IUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := make(object.H)
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		id, err := u.CreateOne(data)
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

// getAPIStore find a store record
//
// GET /api/store/:id
//
// query:
//   fields ([]string) fields to return, id is always returned
//
// success response:
//   {
//      "data": {
//        "store": {
//           "name": "", // store name
//           "description": "", // store description
//           "is_open": bool, // is store open, according to business_hours
//           "average_stars": float, // average stars of all comments
//           "business_hours": [ // a list of business hours from monday to sunday, day without
//                               // any rule is passed null
//              [ // a list of business hours of the day
//                 [ "00:00", "23:59" ], // business hours of the day
//              ]
//           ],
//           "menu": [] // refer to postAPIUsers, price field is not passed if variants exist
//        } // a store document
//    	}
//   }
//
func getAPIUser(u user.IUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := u.FindOneByID(c.Param("id"), c.Query("fields"))
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"user": result,
			},
		})
	}
}

// putAPIStore update a store record
//
// PUT /api/Users
//
// body: (json)
//   {
//      // passed fields are updated (with non-zero value)
//      "name": "", // store display name, max length is 50
//      "description": "", // store description, max length is 1000
//      "business_hours": [] // refer to postAPIUsers, store business hours
//      "categories": [], // []string, categories of the store
//      "price_level": ""// must be one of the values: "cheap", "medium", "expensive"
//      "menu": [] // refer to postAPIUsers, menu does not support operators, when updating,
//                 // pass the entire menu
//      ]
//   }
//
// success response:
//   {
//      "id": "" // store id
//   }
//
func putAPIUser(u user.IUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := make(object.H)
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		data["id"] = c.Param("id")
		if err := u.UpdateOne(data); err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
	}
}

// deleteAPIStore delete a store record
// this endpoint require administrator privilege
//
// DELETE /api/Users
//
// success response:
//   {
//      "deleted_id": "" // store id
//   }
//
func deleteAPIUser(u user.IUserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := u.DeleteOne(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"deleted_id": c.Param("id"),
		})
	}
}
