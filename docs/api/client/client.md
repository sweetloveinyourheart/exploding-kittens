# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [clientserver.proto](#clientserver-proto)
    - [Card](#com-sweetloveinyourheart-kittens-clients-Card)
    - [CreateLobbyRequest](#com-sweetloveinyourheart-kittens-clients-CreateLobbyRequest)
    - [CreateLobbyResponse](#com-sweetloveinyourheart-kittens-clients-CreateLobbyResponse)
    - [CreateNewGuestUserRequest](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserRequest)
    - [CreateNewGuestUserResponse](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserResponse)
    - [DefuseExplodingKittenRequest](#com-sweetloveinyourheart-kittens-clients-DefuseExplodingKittenRequest)
    - [DrawCardRequest](#com-sweetloveinyourheart-kittens-clients-DrawCardRequest)
    - [Game](#com-sweetloveinyourheart-kittens-clients-Game)
    - [Game.Desk](#com-sweetloveinyourheart-kittens-clients-Game-Desk)
    - [Game.Player](#com-sweetloveinyourheart-kittens-clients-Game-Player)
    - [Game.PlayerHand](#com-sweetloveinyourheart-kittens-clients-Game-PlayerHand)
    - [Game.PlayerHandsEntry](#com-sweetloveinyourheart-kittens-clients-Game-PlayerHandsEntry)
    - [GameMetaData](#com-sweetloveinyourheart-kittens-clients-GameMetaData)
    - [GetGameMetaDataRequest](#com-sweetloveinyourheart-kittens-clients-GetGameMetaDataRequest)
    - [GetGameMetaDataResponse](#com-sweetloveinyourheart-kittens-clients-GetGameMetaDataResponse)
    - [GetLobbyReply](#com-sweetloveinyourheart-kittens-clients-GetLobbyReply)
    - [GetLobbyRequest](#com-sweetloveinyourheart-kittens-clients-GetLobbyRequest)
    - [GiveCardRequest](#com-sweetloveinyourheart-kittens-clients-GiveCardRequest)
    - [GuestLoginRequest](#com-sweetloveinyourheart-kittens-clients-GuestLoginRequest)
    - [GuestLoginResponse](#com-sweetloveinyourheart-kittens-clients-GuestLoginResponse)
    - [JoinLobbyRequest](#com-sweetloveinyourheart-kittens-clients-JoinLobbyRequest)
    - [JoinLobbyResponse](#com-sweetloveinyourheart-kittens-clients-JoinLobbyResponse)
    - [LeaveLobbyRequest](#com-sweetloveinyourheart-kittens-clients-LeaveLobbyRequest)
    - [LeaveLobbyResponse](#com-sweetloveinyourheart-kittens-clients-LeaveLobbyResponse)
    - [Lobby](#com-sweetloveinyourheart-kittens-clients-Lobby)
    - [PeekCardsRequest](#com-sweetloveinyourheart-kittens-clients-PeekCardsRequest)
    - [PeekCardsResponse](#com-sweetloveinyourheart-kittens-clients-PeekCardsResponse)
    - [PlayCardsRequest](#com-sweetloveinyourheart-kittens-clients-PlayCardsRequest)
    - [PlayersProfileRequest](#com-sweetloveinyourheart-kittens-clients-PlayersProfileRequest)
    - [PlayersProfileResponse](#com-sweetloveinyourheart-kittens-clients-PlayersProfileResponse)
    - [RetrieveCardsDataResponse](#com-sweetloveinyourheart-kittens-clients-RetrieveCardsDataResponse)
    - [SelectAffectedPlayerRequest](#com-sweetloveinyourheart-kittens-clients-SelectAffectedPlayerRequest)
    - [StartMatchRequest](#com-sweetloveinyourheart-kittens-clients-StartMatchRequest)
    - [StealCardRequest](#com-sweetloveinyourheart-kittens-clients-StealCardRequest)
    - [StreamGameReply](#com-sweetloveinyourheart-kittens-clients-StreamGameReply)
    - [StreamGameRequest](#com-sweetloveinyourheart-kittens-clients-StreamGameRequest)
    - [User](#com-sweetloveinyourheart-kittens-clients-User)
    - [UserProfileResponse](#com-sweetloveinyourheart-kittens-clients-UserProfileResponse)
  
    - [Game.Phase](#com-sweetloveinyourheart-kittens-clients-Game-Phase)
  
    - [ClientServer](#com-sweetloveinyourheart-kittens-clients-ClientServer)
  
- [Scalar Value Types](#scalar-value-types)



<a name="clientserver-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## clientserver.proto



<a name="com-sweetloveinyourheart-kittens-clients-Card"></a>

### Card



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| card_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| code | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-CreateLobbyRequest"></a>

### CreateLobbyRequest
Message for create a lobby


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






<a name="com-sweetloveinyourheart-kittens-clients-DefuseExplodingKittenRequest"></a>

### DefuseExplodingKittenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| card_id | [string](#string) | optional |  |






<a name="com-sweetloveinyourheart-kittens-clients-DrawCardRequest"></a>

### DrawCardRequest
Message for drawing cards
This message is used to draw cards from the deck


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-Game"></a>

### Game



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| game_phase | [Game.Phase](#com-sweetloveinyourheart-kittens-clients-Game-Phase) |  |  |
| player_turn | [string](#string) |  |  |
| players | [Game.Player](#com-sweetloveinyourheart-kittens-clients-Game-Player) | repeated |  |
| player_hands | [Game.PlayerHandsEntry](#com-sweetloveinyourheart-kittens-clients-Game-PlayerHandsEntry) | repeated |  |
| desk | [Game.Desk](#com-sweetloveinyourheart-kittens-clients-Game-Desk) |  |  |
| executing_action | [string](#string) |  |  |
| affected_player | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-Game-Desk"></a>

### Game.Desk



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| desk_id | [string](#string) |  |  |
| remaining_cards | [int32](#int32) |  |  |
| discard_pile | [string](#string) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-Game-Player"></a>

### Game.Player



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player_id | [string](#string) |  |  |
| active | [bool](#bool) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-Game-PlayerHand"></a>

### Game.PlayerHand



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| remaining_cards | [int32](#int32) |  |  |
| hands | [string](#string) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-Game-PlayerHandsEntry"></a>

### Game.PlayerHandsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Game.PlayerHand](#com-sweetloveinyourheart-kittens-clients-Game-PlayerHand) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-GameMetaData"></a>

### GameMetaData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| players | [string](#string) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-GetGameMetaDataRequest"></a>

### GetGameMetaDataRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-GetGameMetaDataResponse"></a>

### GetGameMetaDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| meta | [GameMetaData](#com-sweetloveinyourheart-kittens-clients-GameMetaData) |  |  |






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






<a name="com-sweetloveinyourheart-kittens-clients-GiveCardRequest"></a>

### GiveCardRequest
Message for giving a card
This message is used to give a card to another player


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| card_id | [string](#string) |  |  |






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
| user | [User](#com-sweetloveinyourheart-kittens-clients-User) |  | The user basic info |
| token | [string](#string) |  | The session token for this user. |






<a name="com-sweetloveinyourheart-kittens-clients-JoinLobbyRequest"></a>

### JoinLobbyRequest
Message for join a lobby


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_code | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-JoinLobbyResponse"></a>

### JoinLobbyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-LeaveLobbyRequest"></a>

### LeaveLobbyRequest
Message for leave a lobby


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-LeaveLobbyResponse"></a>

### LeaveLobbyResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-Lobby"></a>

### Lobby



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |
| lobby_code | [string](#string) |  |  |
| lobby_name | [string](#string) |  |  |
| host_user_id | [string](#string) |  |  |
| participants | [string](#string) | repeated |  |
| match_id | [string](#string) | optional |  |






<a name="com-sweetloveinyourheart-kittens-clients-PeekCardsRequest"></a>

### PeekCardsRequest
Message for peeking cards
This message is used to peek at the top card of the deck


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| desk_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-PeekCardsResponse"></a>

### PeekCardsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| card_ids | [string](#string) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-PlayCardsRequest"></a>

### PlayCardsRequest
Message for playing cards
This message is used to play cards in the game


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| card_ids | [string](#string) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-PlayersProfileRequest"></a>

### PlayersProfileRequest
Message for players profile


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_ids | [string](#string) | repeated | Required: UUID of the guest user |






<a name="com-sweetloveinyourheart-kittens-clients-PlayersProfileResponse"></a>

### PlayersProfileResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users | [User](#com-sweetloveinyourheart-kittens-clients-User) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-RetrieveCardsDataResponse"></a>

### RetrieveCardsDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cards | [Card](#com-sweetloveinyourheart-kittens-clients-Card) | repeated |  |






<a name="com-sweetloveinyourheart-kittens-clients-SelectAffectedPlayerRequest"></a>

### SelectAffectedPlayerRequest
Message for selecting affected players
This message is used to select affected players in the game


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| player_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-StartMatchRequest"></a>

### StartMatchRequest
Message for start a match


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lobby_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-StealCardRequest"></a>

### StealCardRequest
Message for stealing a card
This message is used to steal a card from another player


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| card_id | [string](#string) | optional |  |
| card_index | [int32](#int32) | optional |  |






<a name="com-sweetloveinyourheart-kittens-clients-StreamGameReply"></a>

### StreamGameReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_state | [Game](#com-sweetloveinyourheart-kittens-clients-Game) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-StreamGameRequest"></a>

### StreamGameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |
| username | [string](#string) |  |  |
| full_name | [string](#string) |  |  |
| status | [int32](#int32) |  |  |






<a name="com-sweetloveinyourheart-kittens-clients-UserProfileResponse"></a>

### UserProfileResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-kittens-clients-User) |  |  |





 


<a name="com-sweetloveinyourheart-kittens-clients-Game-Phase"></a>

### Game.Phase


| Name | Number | Description |
| ---- | ------ | ----------- |
| INITIALIZING | 0 | Setting up players, shuffling and dealing cards, inserting Exploding Kittens and Defuse cards into the deck |
| TURN_START | 1 | Active player begins their turn |
| ACTION_PHASE | 2 | Player can play as many action cards as they want |
| CARD_DRAWING | 3 | Player draws one card from the deck (mandatory if they didn&#39;t Skip/Attack) |
| TURN_END | 4 | Finalize the turn, next player becomes active |
| GAME_OVER | 5 | When only one player remains |


 

 


<a name="com-sweetloveinyourheart-kittens-clients-ClientServer"></a>

### ClientServer


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RetrieveCardsData | [.google.protobuf.Empty](#google-protobuf-Empty) | [RetrieveCardsDataResponse](#com-sweetloveinyourheart-kittens-clients-RetrieveCardsDataResponse) |  |
| CreateNewGuestUser | [CreateNewGuestUserRequest](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserRequest) | [CreateNewGuestUserResponse](#com-sweetloveinyourheart-kittens-clients-CreateNewGuestUserResponse) |  |
| GuestLogin | [GuestLoginRequest](#com-sweetloveinyourheart-kittens-clients-GuestLoginRequest) | [GuestLoginResponse](#com-sweetloveinyourheart-kittens-clients-GuestLoginResponse) |  |
| GetUserProfile | [.google.protobuf.Empty](#google-protobuf-Empty) | [UserProfileResponse](#com-sweetloveinyourheart-kittens-clients-UserProfileResponse) |  |
| GetPlayersProfile | [PlayersProfileRequest](#com-sweetloveinyourheart-kittens-clients-PlayersProfileRequest) | [PlayersProfileResponse](#com-sweetloveinyourheart-kittens-clients-PlayersProfileResponse) |  |
| CreateLobby | [CreateLobbyRequest](#com-sweetloveinyourheart-kittens-clients-CreateLobbyRequest) | [CreateLobbyResponse](#com-sweetloveinyourheart-kittens-clients-CreateLobbyResponse) |  |
| GetLobby | [GetLobbyRequest](#com-sweetloveinyourheart-kittens-clients-GetLobbyRequest) | [GetLobbyReply](#com-sweetloveinyourheart-kittens-clients-GetLobbyReply) |  |
| StreamLobby | [GetLobbyRequest](#com-sweetloveinyourheart-kittens-clients-GetLobbyRequest) | [GetLobbyReply](#com-sweetloveinyourheart-kittens-clients-GetLobbyReply) stream |  |
| JoinLobby | [JoinLobbyRequest](#com-sweetloveinyourheart-kittens-clients-JoinLobbyRequest) | [JoinLobbyResponse](#com-sweetloveinyourheart-kittens-clients-JoinLobbyResponse) |  |
| LeaveLobby | [LeaveLobbyRequest](#com-sweetloveinyourheart-kittens-clients-LeaveLobbyRequest) | [LeaveLobbyResponse](#com-sweetloveinyourheart-kittens-clients-LeaveLobbyResponse) |  |
| StartMatch | [StartMatchRequest](#com-sweetloveinyourheart-kittens-clients-StartMatchRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| GetGameMetaData | [GetGameMetaDataRequest](#com-sweetloveinyourheart-kittens-clients-GetGameMetaDataRequest) | [GetGameMetaDataResponse](#com-sweetloveinyourheart-kittens-clients-GetGameMetaDataResponse) |  |
| StreamGame | [StreamGameRequest](#com-sweetloveinyourheart-kittens-clients-StreamGameRequest) | [StreamGameReply](#com-sweetloveinyourheart-kittens-clients-StreamGameReply) stream |  |
| PlayCards | [PlayCardsRequest](#com-sweetloveinyourheart-kittens-clients-PlayCardsRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| PeekCards | [PeekCardsRequest](#com-sweetloveinyourheart-kittens-clients-PeekCardsRequest) | [PeekCardsResponse](#com-sweetloveinyourheart-kittens-clients-PeekCardsResponse) |  |
| DrawCard | [DrawCardRequest](#com-sweetloveinyourheart-kittens-clients-DrawCardRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| SelectAffectedPlayer | [SelectAffectedPlayerRequest](#com-sweetloveinyourheart-kittens-clients-SelectAffectedPlayerRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| StealCard | [StealCardRequest](#com-sweetloveinyourheart-kittens-clients-StealCardRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| GiveCard | [GiveCardRequest](#com-sweetloveinyourheart-kittens-clients-GiveCardRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| DefuseExplodingKitten | [DefuseExplodingKittenRequest](#com-sweetloveinyourheart-kittens-clients-DefuseExplodingKittenRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |

 



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

