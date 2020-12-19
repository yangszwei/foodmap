package store

import (
	"foodmap/internal/entity/store"
	"foodmap/internal/infra/delivery"
	"foodmap/internal/infra/object"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupAPI add api routes to server
func SetupAPI(s *delivery.Server, u store.IStoreUsecase) {
	stores := s.API.Group("/stores")
	stores.GET("/", getAPIStores(u))
	stores.POST("/", postAPIStores(u))
	stores.GET("/:id", getAPIStore(u))
	stores.PUT("/:id", putAPIStore(u))
	stores.DELETE("/:id", deleteAPIStore(u))
	stores.GET("/:id/comments", getAPIComments(u))
	stores.POST("/:id/comments", postAPIComments(u))
	stores.DELETE("/:id/comments/:cid", deleteAPIComment(u))
}

// getAPIStores find a list of stores
//
// GET /api/stores
//
// query:
//   query (string) text to search with
//   categories ([]string) filter: only stores with all requested categories
//     are returned
//   limit (int64) limit the number of returned records, ignored when limit and
//     skip are both 0
//   skip (int64) records to skip from search results
//   fields ([]string) fields to return, id is always returned
//
// success response:
//   {
//      "data": {
//         "stores": [
//            {} // a store document, refer to getAPIStore
//         ]
//    	}
//   }
//
func getAPIStores(u store.IStoreUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, s, err := delivery.ParseLimitAndSkip(c.Query("limit"), c.Query("skip"))
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
		}
		result, err := u.Find(c.Query("query"), c.Query("categories"), c.Query("fields"), l, s)
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"stores": result,
			},
		})
	}
}

// postAPIStores create a store record
//
// POST /api/stores
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
func postAPIStores(u store.IStoreUsecase) gin.HandlerFunc {
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
//           "menu": [] // refer to postAPIStores, price field is not passed if variants exist
//        } // a store document
//    	}
//   }
//
func getAPIStore(u store.IStoreUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := u.FindOneByID(c.Param("id"), c.Query("fields"))
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"store": result,
			},
		})
	}
}

// putAPIStore update a store record
//
// PUT /api/stores
//
// body: (json)
//   {
//      // passed fields are updated (with non-zero value)
//      "name": "", // store display name, max length is 50
//      "description": "", // store description, max length is 1000
//      "business_hours": [] // refer to postAPIStores, store business hours
//      "categories": [], // []string, categories of the store
//      "price_level": ""// must be one of the values: "cheap", "medium", "expensive"
//      "menu": [] // refer to postAPIStores, menu does not support operators, when updating,
//                 // pass the entire menu
//      ]
//   }
//
// success response:
//   {
//      "id": "" // store id
//   }
//
func putAPIStore(u store.IStoreUsecase) gin.HandlerFunc {
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
// DELETE /api/stores
//
// success response:
//   {
//      "deleted_id": "" // store id
//   }
//
func deleteAPIStore(u store.IStoreUsecase) gin.HandlerFunc {
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

// getAPIStores find a list of comments
//
// GET /api/stores/:id/comments
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
//         "store": {
//            "comments": [
//               {
//                  "user_id": "", // user id
//                  "user_name": "", // user name,
//                  "stars": 0, // stars
//                  "message": "", // message
//                  "ip_addr": "", // ip address, this field require administrator privilege
//                  "user_agent": "", // user agent, this field require administrator privilege
//               }
//            ]
//         }
//    	}
//   }
//
func getAPIComments(u store.IStoreUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, s, err := delivery.ParseLimitAndSkip(c.Query("limit"), c.Query("skip"))
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		result, err := u.FindComments(c.Param("id"), false, l, s)
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"store": gin.H{
					"comments": result,
				},
			},
		})
	}
}

// postAPIComments add a comment record to store
//
// POST /api/stores/:id/comments
//
// body: (json)
//   {
//      "user_id": "", // required, user id
//      "stars": 0, // required, must be between 0 ~ 5
//      "message": "", // optional, max length is 500
//      "ip_addr": "" // required, client ip address
//      "user_agent": "" // required, client user agent
//   }
//
// success response:
//   {
//     "id": "" // comment id
//   }
//
func postAPIComments(u store.IStoreUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := make(object.H)
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		id, err := u.CreateComment(c.Param("id"), data)
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}

// deleteAPIComment remove a comment from store
// this endpoint require administrator privilege
//
// DELETE /api/stores/:store_id/comments/:comment_id
//
// success response:
//   {
//      "deleted_id": "" // store id
//   }
//
func deleteAPIComment(u store.IStoreUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := u.DeleteComment(c.Param("id"), c.Param("cid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.Error(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"deleted_id": c.Param("cid"),
		})
	}
}
