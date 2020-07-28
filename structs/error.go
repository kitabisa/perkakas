package structs

import (
	"net/http"
)

var ErrUnknown *ErrorResponse = &ErrorResponse{
	Response: Response{
		ResponseCode: "00001",
		ResponseDesc: ResponseDesc{
			ID: "Ups ada kesalahan, silahkan coba beberapa saat lagi",
			EN: "Unknown error",
		},
	},
	HttpStatus: http.StatusInternalServerError,
}

var ErrUnauthorized *ErrorResponse = &ErrorResponse{
	Response: Response{
		ResponseCode: "00002",
		ResponseDesc: ResponseDesc{
			ID: "Anda tidak diijinkan",
			EN: "You are not authorized",
		},
	},
	HttpStatus: http.StatusUnauthorized,
}

var ErrInvalidHeader *ErrorResponse = &ErrorResponse{
	Response: Response{
		ResponseCode: "00003",
		ResponseDesc: ResponseDesc{
			ID: "Header tidak valid atau tidak lengkap",
			EN: "Invalid/incomplete header",
		},
	},
	HttpStatus: http.StatusBadRequest,
}

var ErrInvalidHeaderSignature *ErrorResponse = &ErrorResponse{
	Response: Response{
		ResponseCode: "00004",
		ResponseDesc: ResponseDesc{
			ID: "Header signature tidak valid",
			EN: "Invalid header signature",
		},
	},
	HttpStatus: http.StatusBadRequest,
}

var ErrInvalidHeaderTime *ErrorResponse = &ErrorResponse{
	Response: Response{
		ResponseCode: "00005",
		ResponseDesc: ResponseDesc{
			ID: "Request sudah kedaluwarsa",
			EN: "Request already expired",
		},
	},
	HttpStatus: http.StatusBadRequest,
}
