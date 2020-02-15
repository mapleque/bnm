package server

// order stage change process
//
//  customer order       customer cancel
//  ---------------> new ---------------> c_cancel
//                    |
//                    |  business cancel
//                    + ----------------> b_cancel
//                    |
//       customer pay |
//                    |
//                    v   customer want cancel
//                   paid --------------------> c_w_cancel
//                    |                            |
//                    |                            |
//                    +----------------------------+
//                    |                            |
//     business trans |            business repaid |
//                    |                            |
//                    v                            v
//                 transport                     repaid
//                    |
//                    |
//    customer commit |
//                    |
//                    v
//                  finish

const (
	// status = 0
	ORDER_STAGE_NEW        = "new"
	ORDER_STAGE_PAID       = "paid"
	ORDER_STAGE_TRANSPORT  = "transport"
	ORDER_STAGE_C_W_CANCEL = "c_w_cancel"

	// status = 1
	ORDER_STAGE_FINISH   = "finish"
	ORDER_STAGE_C_CANCEL = "c_cancel"
	ORDER_STAGE_B_CANCEL = "b_cancel"
	ORDER_STAGE_REPAID   = "repaid"
)
