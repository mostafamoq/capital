package capital

import "net/textproto"

var (
	HttpContentTypeHeaderKey        = textproto.CanonicalMIMEHeaderKey("Content-Type")
	HttpContentTypeHeaderJson       = textproto.CanonicalMIMEHeaderKey("application/json")
	WSProtocolHeaderKey             = textproto.CanonicalMIMEHeaderKey("Sec-WebSocket-Protocol")
	HttpAuthorizationHeaderKey      = textproto.CanonicalMIMEHeaderKey("Authorization")
	HttpAuthorizationHeaderTokenKey = textproto.CanonicalMIMEHeaderKey("X-Auth-Token")
	HttpAcceptHeaderKey             = textproto.CanonicalMIMEHeaderKey("Accept")
	HttpMergePatchStrategic         = textproto.CanonicalMIMEHeaderKey("application/strategic-merge-patch+json")
	HttpMergePatchJsonStrategic     = textproto.CanonicalMIMEHeaderKey("application/merge-patch+json")
	HttpJsonPatchStrategic          = textproto.CanonicalMIMEHeaderKey("application/json-patch+json")
)
