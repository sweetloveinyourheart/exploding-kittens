# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [dataprovider.proto](#dataprovider-proto)
    - [Card](#com-sweetloveinyourheart-kittens-dataproviders-Card)
    - [GetCardsResponse](#com-sweetloveinyourheart-kittens-dataproviders-GetCardsResponse)
    - [GetMapCardsResponse](#com-sweetloveinyourheart-kittens-dataproviders-GetMapCardsResponse)
    - [GetMapCardsResponse.CardsEntry](#com-sweetloveinyourheart-kittens-dataproviders-GetMapCardsResponse-CardsEntry)
  
    - [DataProvider](#com-sweetloveinyourheart-kittens-dataproviders-DataProvider)
  
- [Scalar Value Types](#scalar-value-types)



<a name="dataprovider-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## dataprovider.proto



<a name="com-sweetloveinyourheart-kittens-dataproviders-Card"></a>

### Card
The Card message definition


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| card_id | [string](#string) |  |  |
| code | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| quantity | [int32](#int32) |  |  |
| effects | [bytes](#bytes) |  |  |
| combo_effects | [bytes](#bytes) |  |  |






<a name="com-sweetloveinyourheart-kittens-dataproviders-GetCardsResponse"></a>

### GetCardsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cards | [Card](#com-sweetloveinyourheart-kittens-dataproviders-Card) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-dataproviders-GetMapCardsResponse"></a>

### GetMapCardsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cards | [GetMapCardsResponse.CardsEntry](#com-sweetloveinyourheart-kittens-dataproviders-GetMapCardsResponse-CardsEntry) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-dataproviders-GetMapCardsResponse-CardsEntry"></a>

### GetMapCardsResponse.CardsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Card](#com-sweetloveinyourheart-kittens-dataproviders-Card) |  |  |





 

 

 


<a name="com-sweetloveinyourheart-kittens-dataproviders-DataProvider"></a>

### DataProvider


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetCards | [.google.protobuf.Empty](#google-protobuf-Empty) | [GetCardsResponse](#com-sweetloveinyourheart-kittens-dataproviders-GetCardsResponse) | Get cards |
| GetMapCards | [.google.protobuf.Empty](#google-protobuf-Empty) | [GetMapCardsResponse](#com-sweetloveinyourheart-kittens-dataproviders-GetMapCardsResponse) | Get cards as map |

 



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

