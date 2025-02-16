# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [clientserver.proto](#clientserver-proto)
    - [CreateLobbyRequest](#com-sweetloveinyourheart-kittens-clients-CreateLobbyRequest)
    - [CreateLobbyResponse](#com-sweetloveinyourheart-kittens-clients-CreateLobbyResponse)
    - [CreateNewGuestUserRequest](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserRequest)
    - [CreateNewGuestUserResponse](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserResponse)
    - [GetLobbyReply](#com-sweetloveinyourheart-kittens-clients-GetLobbyReply)
    - [GetLobbyRequest](#com-sweetloveinyourheart-kittens-clients-GetLobbyRequest)
    - [GuestLoginRequest](#com-sweetloveinyourheart-kittens-clients-GuestLoginRequest)
    - [GuestLoginResponse](#com-sweetloveinyourheart-kittens-clients-GuestLoginResponse)
    - [Lobby](#com-sweetloveinyourheart-kittens-clients-Lobby)
    - [PlayerProfileResponse](#com-sweetloveinyourheart-kittens-clients-PlayerProfileResponse)
    - [User](#com-sweetloveinyourheart-kittens-clients-User)
  
    - [ClientServer](#com-sweetloveinyourheart-kittens-clients-ClientServer)
  
- [Scalar Value Types](#scalar-value-types)



<a name="clientserver-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## clientserver.proto



<a name="com-sweetloveinyourheart-kittens-clients-CreateLobbyRequest"></a>

### CreateLobbyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_name | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-CreateLobbyResponse"></a>

### CreateLobbyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserRequest"></a>

### CreateNewGuestUserRequest
Message for creating a new guest user


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| username | [string](#string) |  | Required: Username of the guest user |
| full_name | [string](#string) |  | Required: Full name of the guest user |






<a name="com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserResponse"></a>

### CreateNewGuestUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-kittens-clients-User) |  | The user basic info |






<a name="com-sweetloveinyourheart-kittens-clients-GetLobbyReply"></a>

### GetLobbyReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby | [Lobby](#com-sweetloveinyourheart-kittens-clients-Lobby) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-GetLobbyRequest"></a>

### GetLobbyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-GuestLoginRequest"></a>

### GuestLoginRequest
Message for guest login


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  | Required: UUID of the guest user |






<a name="com-sweetloveinyourheart-kittens-clients-GuestLoginResponse"></a>

### GuestLoginResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  | The database id for this user (UUID). |
| token | [string](#string) |  | The session token for this user. |






<a name="com-sweetloveinyourheart-kittens-clients-Lobby"></a>

### Lobby



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |
| lobby_code | [string](#string) |  |  |
| lobby_name | [string](#string) |  |  |
| host_user_id | [string](#string) |  |  |
| participants | [string](#string) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-PlayerProfileResponse"></a>

### PlayerProfileResponse
Message for player profile


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-kittens-clients-User) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |
| username | [string](#string) |  |  |
| full_name | [string](#string) |  |  |
| status | [int32](#int32) |  |  |





 

 

 


<a name="com-sweetloveinyourheart-kittens-clients-ClientServer"></a>

### ClientServer


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateNewGuestUser | [CreateNewGuestUserRequest](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserRequest) | [CreateNewGuestUserResponse](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserResponse) |  |
| GuestLogin | [GuestLoginRequest](#com-sweetloveinyourheart-kittens-clients-GuestLoginRequest) | [GuestLoginResponse](#com-sweetloveinyourheart-kittens-clients-GuestLoginResponse) |  |
| GetPlayerProfile | [.google.protobuf.Empty](#google-protobuf-Empty) | [PlayerProfileResponse](#com-sweetloveinyourheart-kittens-clients-PlayerProfileResponse) |  |
| CreateLobby | [CreateLobbyRequest](#com-sweetloveinyourheart-kittens-clients-CreateLobbyRequest) | [CreateLobbyResponse](#com-sweetloveinyourheart-kittens-clients-CreateLobbyResponse) |  |
| StreamLobby | [GetLobbyRequest](#com-sweetloveinyourheart-kittens-clients-GetLobbyRequest) | [GetLobbyReply](#com-sweetloveinyourheart-kittens-clients-GetLobbyReply) stream |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

