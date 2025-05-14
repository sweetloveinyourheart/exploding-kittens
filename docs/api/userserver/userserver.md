# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [userserver.proto](#userserver-proto)
    - [CreateUserRequest](#com-sweetloveinyourheart-kittens-users-CreateUserRequest)
    - [CreateUserResponse](#com-sweetloveinyourheart-kittens-users-CreateUserResponse)
    - [GetUserRequest](#com-sweetloveinyourheart-kittens-users-GetUserRequest)
    - [GetUserResponse](#com-sweetloveinyourheart-kittens-users-GetUserResponse)
    - [SignInRequest](#com-sweetloveinyourheart-kittens-users-SignInRequest)
    - [SignInResponse](#com-sweetloveinyourheart-kittens-users-SignInResponse)
    - [User](#com-sweetloveinyourheart-kittens-users-User)
  
    - [CreateUserRequest.AuthProvider](#com-sweetloveinyourheart-kittens-users-CreateUserRequest-AuthProvider)
  
    - [UserServer](#com-sweetloveinyourheart-kittens-users-UserServer)
  
- [Scalar Value Types](#scalar-value-types)



<a name="userserver-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## userserver.proto



<a name="com-sweetloveinyourheart-kittens-users-CreateUserRequest"></a>

### CreateUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| username | [string](#string) |  |  |
| full_name | [string](#string) |  |  |
| auth_provider | [CreateUserRequest.AuthProvider](#com-sweetloveinyourheart-kittens-users-CreateUserRequest-AuthProvider) |  |  |
| meta | [string](#string) | optional |  |






<a name="com-sweetloveinyourheart-kittens-users-CreateUserResponse"></a>

### CreateUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-kittens-users-User) |  |  |






<a name="com-sweetloveinyourheart-kittens-users-GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-users-GetUserResponse"></a>

### GetUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-kittens-users-User) |  |  |






<a name="com-sweetloveinyourheart-kittens-users-SignInRequest"></a>

### SignInRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-users-SignInResponse"></a>

### SignInResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-kittens-users-User) |  | The user basic info |
| token | [string](#string) |  | The session token for this user. |






<a name="com-sweetloveinyourheart-kittens-users-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |
| username | [string](#string) |  |  |
| full_name | [string](#string) |  |  |
| status | [int32](#int32) |  |  |
| created_at | [int64](#int64) |  | Unix time for CreatedAt |
| updated_at | [int64](#int64) |  | Unix time for UpdatedAt |





 


<a name="com-sweetloveinyourheart-kittens-users-CreateUserRequest-AuthProvider"></a>

### CreateUserRequest.AuthProvider


| Name | Number | Description |
| ---- | ------ | ----------- |
| GUEST | 0 | Guest user |
| GOOGLE | 1 | Google SSO user |


 

 


<a name="com-sweetloveinyourheart-kittens-users-UserServer"></a>

### UserServer


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetUser | [GetUserRequest](#com-sweetloveinyourheart-kittens-users-GetUserRequest) | [GetUserResponse](#com-sweetloveinyourheart-kittens-users-GetUserResponse) | Get a user by user_id |
| CreateNewUser | [CreateUserRequest](#com-sweetloveinyourheart-kittens-users-CreateUserRequest) | [CreateUserResponse](#com-sweetloveinyourheart-kittens-users-CreateUserResponse) | Create new user |
| SignIn | [SignInRequest](#com-sweetloveinyourheart-kittens-users-SignInRequest) | [SignInResponse](#com-sweetloveinyourheart-kittens-users-SignInResponse) | Sign in |

 



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

