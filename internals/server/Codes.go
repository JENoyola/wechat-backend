package server

// INFORMATION CODES (mostly internal server use)

// 2xx SUCCESS CODES
/*
	OK

	mean all the operations where ok.
	used whe no data is been returned or client is not handling.
	any client redirection.
*/
const OK = 206

/*
COMPLETED

means all the operations where completed.
mainly used when the operation is one way - Client -> Server.
and used when no data is returned.
*/
const COMPLETED = 299

// 3xx redirection
/*
	REDIRECTION

	means all operations where successfull and the user needs
	to be redirected to a screen. Most of the time returns data.
*/
const REDIRECTION = 306

// 4xx CLIENT ERRORS

/*
BAD_REQUEST

means the request was bad in the context of not including
essencial data.
This error is mainly used in json payloads
*/
const BAD_REQUEST = 400

/*
BAD_CREDENTIALS

means the payload that was sent, contains invalid credentials
and user is not allow to proceed.
This code is mainly used for password checks.
*/
const BAD_CREDENTIALS = 402

/*
NO_DOCUMENTS

means the information given didnt get any results in the database
because no registry was found.
This code is mainly used for database retrieve requests
*/
const NO_DOCUMENTS = 404

/*
DOCUMENT_FOUND

means that the information retrieved a document in the registry.
This code is mainly used when it needs to check for availability
*/
const DOCUMENT_FOUND = 404

/*
BAD_FIELD

means that the user introduced an invalid information in a file of a
request.
This code is mainly used for json requests.
Ej. an empty string, an uninitialized field, etc.
*/
const BAD_FIELD = 408

/*
NOT_ALLOWED

means that the user does not have the right credentials for certain operations.
This code is mainly used for changes that the user is requsting to do.
*/
const NOT_ALLOWED = 411

// 5xx SERVER ERRORS

/*
SERVER_ERROR

means that an internal error occured in a function, method, etc.
This code is mainlty used for sending a response to the client explaining that the
server is experiencing errors within its system.
*/
const SERVER_ERROR = 500

/*
DB_ERROR

means that an internal error occurred while calling the database.
This code is specifically and reserved for databases requests.
*/
const DB_ERROR = 505

/*
MAIL_ERROR

means that an internal error occurred while using the mailing service.
This code is specifically and reserved for mail services.
*/
const MAIL_ERROR = 507

/*
SERVICES_ERROR

means that an internal error occured while using third party services.
This code is specifically and reserved for errors that occures while using third party
services.
*/
const SERVICES_ERROR = 511

/*
FOREIGN_ACCESS

means that another device is trying to access the device
*/
const FOREIGN_ACCESS = 525

// WEBSOCKET CODES

/*
	PROVIDER_ERROR
	means that the provider used on websocket throw an error
*/

const PROVIDER_ERROR = 658
