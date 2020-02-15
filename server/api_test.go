package server

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testClient   = &http.Client{}
	testToken    = "900150983cd24fb0d6963f7d28e17f72"
	testBusiness = &Business{
		Id:     1,
		Name:   "for test",
		Avatar: "",
		Desc:   "",
		Qrcode: "",
		Wxid:   "test wx id",
	}
	testItem = &Item{
		Name:  "test item",
		Price: 2,
		Pic:   "test pic",
		Desc:  "test desc",
	}
	testItemUpdate = &Item{
		Price:  1,
		Pic:    "test pic update",
		Desc:   "test desc update",
		Status: 1,
	}
	testAddress = &Address{
		Label:    "test label",
		Reciever: "test reciever",
		Address:  "test address",
		Phone:    "test phone",
	}
	testAddressUpdate = &Address{
		Label:    "test update label",
		Reciever: "test update reciever",
		Address:  "test update address",
		Phone:    "test update phone",
	}
	testOrder = &Order{
		Name:     "test item",
		Price:    1,
		Counts:   2,
		Reciever: "test reciever",
		Address:  "test address",
		Phone:    "test phone",
	}
	testOrderPay = &Order{
		Stage: ORDER_STAGE_PAID,
	}
	testOrderCCancel = &Order{
		Stage: ORDER_STAGE_C_CANCEL,
	}
	testOrderBCancel = &Order{
		Stage: ORDER_STAGE_B_CANCEL,
	}
	testOrderCWCancel = &Order{
		Stage: ORDER_STAGE_C_W_CANCEL,
	}
	testOrderTransport = &Order{
		ExpNo: "test exp no",
		Stage: ORDER_STAGE_TRANSPORT,
	}
	testOrderRepaid = &Order{
		Stage: ORDER_STAGE_REPAID,
	}
	testOrderFinish = &Order{
		Stage: ORDER_STAGE_FINISH,
	}
)

var testCases = []*ApiCase{
	{
		sql:  "INSERT INTO customer_user (id,open_id,token,expired_at) VALUES (1,?,?, DATE_ADD(NOW(),INTERVAL 7 DAY))",
		bind: []interface{}{"testopenid", testToken},
	},
	{
		sql:  "INSERT INTO business_user (id,open_id,token,expired_at) VALUES (1,?,?, DATE_ADD(NOW(),INTERVAL 7 DAY))",
		bind: []interface{}{"testopenid", testToken},
	},
	{
		sql:  "INSERT INTO business_profile (id,buid,name,wxid) VALUES (1,1,?,?)",
		bind: []interface{}{testBusiness.Name, testBusiness.Wxid},
	},
	{msg: "customer address add",
		request:  testReq("POST", "/customer/address", toJson(testAddress)),
		response: testResp(nil),
	},
	{msg: "customer address get list",
		request: testReq("GET", "/customer/address?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Address{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) == 1 &&
				t[0].Reciever == testAddress.Reciever &&
				t[0].Address == testAddress.Address &&
				t[0].Phone == testAddress.Phone &&
				t[0].Label == testAddress.Label
		}),
	},
	{msg: "customer address update",
		request:  testReq("POST", "/customer/address/1", toJson(testAddressUpdate)),
		response: testResp(nil),
	},
	{msg: "customer address get list for check update",
		request: testReq("GET", "/customer/address?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Address{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return t[0].Reciever == testAddressUpdate.Reciever &&
				t[0].Address == testAddressUpdate.Address &&
				t[0].Phone == testAddressUpdate.Phone &&
				t[0].Label == testAddressUpdate.Label
		}),
	},
	{msg: "customer address delete",
		request:  testReq("DELETE", "/customer/address/1", ""),
		response: testResp(nil),
	},
	{msg: "customer address get list for check delete",
		request: testReq("GET", "/customer/address?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Address{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) == 0
		}),
	},
	{msg: "commit a business profile",
		request:  testReq("POST", "/business/profile", toJson(testBusiness)),
		response: testResp(nil),
	},
	{msg: "check business profile status is 0",
		request: testReq("GET", "/business/profile", ""),
		response: testResp(func(resp []byte) bool {
			t := &Business{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Status == 0
		}),
	},
	{msg: "can not get items because is not valid",
		request:  testReq("GET", "/business/items?s=10", ""),
		response: testErrorStatus(STATUS_SERVER_ERROR),
	},
	{msg: "valid profile by sql",
		sql:  "UPDATE business_profile SET status = 1 WHERE id = 1",
		bind: []interface{}{},
	},
	{msg: "profile is valid, status is 1",
		request: testReq("GET", "/business/profile", ""),
		response: testResp(func(resp []byte) bool {
			t := &Business{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Status == 1
		}),
	},
	{msg: "business add a item 1",
		request:  testReq("POST", "/business/items", toJson(testItem)),
		response: testResp(nil),
	},
	{msg: "item list has 1 item",
		request: testReq("GET", "/business/items?s=10", ""),
		response: testResp(func(resp []byte) bool {
			t := []*Item{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) > 0 && t[0].Status == 0 &&
				t[0].Name == testItem.Name &&
				t[0].Price == testItem.Price &&
				t[0].Pic == testItem.Pic &&
				t[0].Desc == testItem.Desc
		}),
	},
	{msg: "customer can not get the item, it's status is 0",
		request: testReq("GET", "/customer/business/1/items?s=10", ""),
		response: testResp(func(resp []byte) bool {
			t := []*Item{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) == 0
		}),
	},
	{msg: "update item status to 1",
		request:  testReq("POST", "/business/items/1", toJson(testItemUpdate)),
		response: testResp(nil),
	},
	{msg: "get item which status is 1",
		request: testReq("GET", "/business/items/1", ""),
		response: testResp(func(resp []byte) bool {
			t := &Item{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Status == testItemUpdate.Status &&
				t.Name != testItemUpdate.Name &&
				t.Price == testItemUpdate.Price &&
				t.Pic == testItemUpdate.Pic &&
				t.Desc == testItemUpdate.Desc
		}),
	},
	{msg: "customer can get business list with 1 record",
		request: testReq("GET", "/customer/business?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Business{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) == 1
		}),
	},
	{msg: "customer can get business 1",
		request: testReq("GET", "/customer/business/1", ``),
		response: testResp(func(resp []byte) bool {
			t := &Business{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Id == 1
		}),
	},
	{msg: "customer can get item list",
		request: testReq("GET", "/customer/business/1/items?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Item{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) == 1
		}),
	},
	{msg: "customer can get item 1",
		request: testReq("GET", "/customer/business/1/items/1", ``),
		response: testResp(func(resp []byte) bool {
			t := &Business{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Id == 1
		}),
	},
	{msg: "customer can make order",
		request:  testReq("POST", "/customer/business/1/items/1", toJson(testOrder)),
		response: testResp(nil),
	},
	{msg: "customer can find order in list",
		request: testReq("GET", "/customer/orders?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Order{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) == 1
		}),
	},
	{msg: "business can find order 1 properties same with test order",
		request: testReq("GET", "/business/orders/1", ""),
		response: testResp(func(resp []byte) bool {
			t := &Order{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Bid == 1 && t.Cuid == 1 &&
				t.Stage == ORDER_STAGE_NEW && t.Status == 0 &&
				t.Name == testOrder.Name &&
				t.Price == testOrder.Price &&
				t.Counts == testOrder.Counts &&
				t.Reciever == testOrder.Reciever &&
				t.Address == testOrder.Address &&
				t.Phone == testOrder.Phone &&
				t.Additional == testOrder.Additional
		}),
	},
	{msg: "customer can find order 1 properties same with test order",
		request: testReq("GET", "/customer/orders/1", ""),
		response: testResp(func(resp []byte) bool {
			t := &Order{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Bid == 1 && t.Cuid == 1 &&
				t.Stage == ORDER_STAGE_NEW && t.Status == 0 &&
				t.Name == testOrder.Name &&
				t.Price == testOrder.Price &&
				t.Counts == testOrder.Counts &&
				t.Reciever == testOrder.Reciever &&
				t.Address == testOrder.Address &&
				t.Phone == testOrder.Phone &&
				t.Additional == testOrder.Additional
		}),
	},
	{msg: "business can find order",
		request: testReq("GET", "/business/orders?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Order{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) == 1
		}),
	},
	{msg: "customer pay for order 1",
		request:  testReq("POST", "/customer/orders/1", toJson(testOrderPay)),
		response: testResp(nil),
	},
	{msg: "business transport order 1",
		request:  testReq("POST", "/business/orders/1", toJson(testOrderTransport)),
		response: testResp(nil),
	},
	{msg: "customer can find order exp no in order 1",
		request: testReq("GET", "/customer/orders/1", ""),
		response: testResp(func(resp []byte) bool {
			t := &Order{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Bid == 1 && t.Cuid == 1 &&
				t.Stage == ORDER_STAGE_TRANSPORT && t.Status == 0 &&
				t.ExpNo == testOrderTransport.ExpNo
		}),
	},
	{msg: "customer finish order 1",
		request:  testReq("POST", "/customer/orders/1", toJson(testOrderFinish)),
		response: testResp(nil),
	},
	{msg: "the order 1 is finished",
		request: testReq("GET", "/customer/orders/1", ""),
		response: testResp(func(resp []byte) bool {
			t := &Order{}
			if err := json.Unmarshal(resp, t); err != nil {
				return false
			}
			return t.Bid == 1 && t.Cuid == 1 &&
				t.Stage == ORDER_STAGE_FINISH && t.Status == 1
		}),
	},
	{msg: "customer make a new order 2",
		request:  testReq("POST", "/customer/business/1/items/1", toJson(testOrder)),
		response: testResp(nil),
	},
	{msg: "customer cancel order 2",
		request:  testReq("POST", "/customer/orders/2", toJson(testOrderCCancel)),
		response: testResp(nil),
	},
	{msg: "customer make a new order 3",
		request:  testReq("POST", "/customer/business/1/items/1", toJson(testOrder)),
		response: testResp(nil),
	},
	{msg: "business cancel order 3",
		request:  testReq("POST", "/business/orders/3", toJson(testOrderBCancel)),
		response: testResp(nil),
	},
	{msg: "customer make a new order 4",
		request:  testReq("POST", "/customer/business/1/items/1", toJson(testOrder)),
		response: testResp(nil),
	},
	{msg: "customer pay order 4",
		request:  testReq("POST", "/customer/orders/4", toJson(testOrderPay)),
		response: testResp(nil),
	},
	{msg: "business repaid order 4",
		request:  testReq("POST", "/business/orders/4", toJson(testOrderRepaid)),
		response: testResp(nil),
	},
	{msg: "customer make a new order 5",
		request:  testReq("POST", "/customer/business/1/items/1", toJson(testOrder)),
		response: testResp(nil),
	},
	{msg: "customer pay order 5",
		request:  testReq("POST", "/customer/orders/5", toJson(testOrderPay)),
		response: testResp(nil),
	},
	{msg: "customer want cancel order 5",
		request:  testReq("POST", "/customer/orders/5", toJson(testOrderCWCancel)),
		response: testResp(nil),
	},
	{msg: "business repaid order 5",
		request:  testReq("POST", "/business/orders/5", toJson(testOrderRepaid)),
		response: testResp(nil),
	},
	{msg: "business can find all 5 orders's status should be 1",
		request: testReq("GET", "/business/orders?s=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Order{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			if len(t) != 5 {
				return false
			}
			for _, o := range t {
				if o.Status != 1 {
					return false
				}
			}
			return true
		}),
	},
	{msg: "order next page is empty",
		request: testReq("GET", "/business/orders?s=10&l=10", ``),
		response: testResp(func(resp []byte) bool {
			t := []*Order{}
			if err := json.Unmarshal(resp, &t); err != nil {
				return false
			}
			return len(t) != 0
		}),
	},
	{msg: "order log numbers is 15",
		query: "SELECT COUNT(*) FROM order_log",
		bind:  []interface{}{},
		queryAssert: func(qr DBQueryResult) bool {
			var c int
			if err := qr.Scan(&c); err != nil {
				return false
			}
			return c == 15
		},
	},
	{msg: "404",
		request:  testReq("GET", "/hello", ""),
		response: testHttpStatus(404),
	},
	{msg: "405",
		request:  testReq("GET", "/business/login", ""),
		response: testHttpStatus(405),
	},
}

type ApiCase struct {
	msg         string
	query       string
	queryAssert func(qr DBQueryResult) bool
	sql         string
	bind        []interface{}
	request     *testRequest
	response    *testResponse
}

func TestApi(t *testing.T) {
	db := NewDB(os.Getenv("DB_DSN"))
	initDBSchema(db, "../sql/")
	s := NewServer(db, nil, nil, 8080)
	ts := httptest.NewServer(s)
	defer ts.Close()

	for _, c := range testCases {
		if c.sql != "" {
			if err := db.Execute(c.sql, c.bind...).Error(); err != nil {
				t.Error(err)
				t.Fatal(c.msg)
			}
		} else if c.request != nil {
			resp, err := testClient.Do(buildRequest(ts.URL, c.request))
			if err != nil {
				t.Error(err)
				t.Fatal(c.msg)
			}
			if !assertResponse(t, resp, c.response) {
				t.Fatal(c.msg)
			}

		} else if c.query != "" {
			qr := db.Query(
				c.query,
				c.bind...,
			)
			if !c.queryAssert(qr) {
				t.Error("query assert faild", c.query)
				t.Fatal(c.msg)
			}
		}
	}

}

func initDBSchema(db DBConn, schemaPath string) {
	sqlDir, err := ioutil.ReadDir(schemaPath)
	if err != nil {
		panic(err)
	}
	for _, fi := range sqlDir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), ".SQL") {
			sqlfile := schemaPath + string(os.PathSeparator) + fi.Name()
			sourceSql, err := ioutil.ReadFile(sqlfile)
			if err != nil {
				panic(err)
			}
			for _, sql := range strings.Split(string(sourceSql), ";") {
				if len(strings.TrimSpace(sql)) > 0 {
					if err := db.Execute(sql).Error(); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

type testRequest struct {
	method string
	url    string
	body   []byte
}

func testReq(method, url, body string) *testRequest {
	var b []byte
	if body != "" {
		b = []byte(body)
	}
	return &testRequest{
		method: method,
		url:    url,
		body:   b,
	}
}

func buildRequest(host string, tr *testRequest) *http.Request {
	var body io.Reader
	if tr.body != nil {
		body = bytes.NewReader(tr.body)
	}
	req, _ := http.NewRequest(tr.method, host+tr.url, body)
	req.Header.Set("SessionKey", testToken)
	return req
}

type testResponse struct {
	statusCode int

	Status  int             `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message interface{}     `json:"message"`

	assertFunc func([]byte) bool
}

func testResp(assertFunc func([]byte) bool) *testResponse {
	return &testResponse{
		statusCode: 200,
		Status:     STATUS_SUCCESS,
		assertFunc: assertFunc,
	}
}

func assertResponse(t *testing.T, resp *http.Response, target *testResponse) bool {
	if target.statusCode > 0 && resp.StatusCode != target.statusCode {
		t.Error("invalid http status code, should be:", target.statusCode, "but:", resp.StatusCode)
		return false
	}

	response := &testResponse{}
	if target.Status > 0 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return false
		}
		if err := json.Unmarshal(body, response); err != nil {
			t.Error("invalid response body", string(body))
			return false
		}
		if response.Status != target.Status {
			t.Error("invalid response status, should be:", target.Status, "but:", response.Status, "with body", string(body))
			return false
		}
	}
	if target.assertFunc != nil && !target.assertFunc(response.Data) {
		t.Error("invalid response data", string(response.Data))
		return false
	}
	return true
}

func toJson(tar interface{}) string {
	bt, _ := json.Marshal(tar)
	return string(bt)
}

func testErrorStatus(status int) *testResponse {
	return &testResponse{
		statusCode: 200,
		Status:     status,
	}
}

func testHttpStatus(status int) *testResponse {
	return &testResponse{
		statusCode: status,
	}
}
