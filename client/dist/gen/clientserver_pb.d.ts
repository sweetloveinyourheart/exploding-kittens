import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.User
 */
export declare class User extends Message<User> {
    /**
     * @generated from field: string user_id = 1;
     */
    userId: string;
    /**
     * @generated from field: string username = 2;
     */
    username: string;
    /**
     * @generated from field: string full_name = 3;
     */
    fullName: string;
    /**
     * @generated from field: int32 status = 4;
     */
    status: number;
    constructor(data?: PartialMessage<User>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.User";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): User;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): User;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): User;
    static equals(a: User | PlainMessage<User> | undefined, b: User | PlainMessage<User> | undefined): boolean;
}
/**
 * Message for creating a new guest user
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserRequest
 */
export declare class CreateNewGuestUserRequest extends Message<CreateNewGuestUserRequest> {
    /**
     * Required: Username of the guest user
     *
     * @generated from field: string username = 1;
     */
    username: string;
    /**
     * Required: Full name of the guest user
     *
     * @generated from field: string full_name = 2;
     */
    fullName: string;
    constructor(data?: PartialMessage<CreateNewGuestUserRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateNewGuestUserRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateNewGuestUserRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateNewGuestUserRequest;
    static equals(a: CreateNewGuestUserRequest | PlainMessage<CreateNewGuestUserRequest> | undefined, b: CreateNewGuestUserRequest | PlainMessage<CreateNewGuestUserRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserResponse
 */
export declare class CreateNewGuestUserResponse extends Message<CreateNewGuestUserResponse> {
    /**
     * The user basic info
     *
     * @generated from field: com.sweetloveinyourheart.kittens.clients.User user = 1;
     */
    user?: User;
    constructor(data?: PartialMessage<CreateNewGuestUserResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateNewGuestUserResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateNewGuestUserResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateNewGuestUserResponse;
    static equals(a: CreateNewGuestUserResponse | PlainMessage<CreateNewGuestUserResponse> | undefined, b: CreateNewGuestUserResponse | PlainMessage<CreateNewGuestUserResponse> | undefined): boolean;
}
/**
 * Message for guest login
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.GuestLoginRequest
 */
export declare class GuestLoginRequest extends Message<GuestLoginRequest> {
    /**
     * Required: UUID of the guest user
     *
     * @generated from field: string user_id = 1;
     */
    userId: string;
    constructor(data?: PartialMessage<GuestLoginRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GuestLoginRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GuestLoginRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GuestLoginRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GuestLoginRequest;
    static equals(a: GuestLoginRequest | PlainMessage<GuestLoginRequest> | undefined, b: GuestLoginRequest | PlainMessage<GuestLoginRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.GuestLoginResponse
 */
export declare class GuestLoginResponse extends Message<GuestLoginResponse> {
    /**
     * The user basic info
     *
     * @generated from field: com.sweetloveinyourheart.kittens.clients.User user = 1;
     */
    user?: User;
    /**
     * The session token for this user.
     *
     * @generated from field: string token = 2;
     */
    token: string;
    constructor(data?: PartialMessage<GuestLoginResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GuestLoginResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GuestLoginResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GuestLoginResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GuestLoginResponse;
    static equals(a: GuestLoginResponse | PlainMessage<GuestLoginResponse> | undefined, b: GuestLoginResponse | PlainMessage<GuestLoginResponse> | undefined): boolean;
}
/**
 * Message for player profile
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.PlayerProfileRequest
 */
export declare class PlayerProfileRequest extends Message<PlayerProfileRequest> {
    /**
     * Required: UUID of the guest user
     *
     * @generated from field: string user_id = 1;
     */
    userId: string;
    constructor(data?: PartialMessage<PlayerProfileRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.PlayerProfileRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PlayerProfileRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PlayerProfileRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PlayerProfileRequest;
    static equals(a: PlayerProfileRequest | PlainMessage<PlayerProfileRequest> | undefined, b: PlayerProfileRequest | PlainMessage<PlayerProfileRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.PlayerProfileResponse
 */
export declare class PlayerProfileResponse extends Message<PlayerProfileResponse> {
    /**
     * @generated from field: com.sweetloveinyourheart.kittens.clients.User user = 1;
     */
    user?: User;
    constructor(data?: PartialMessage<PlayerProfileResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.PlayerProfileResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PlayerProfileResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PlayerProfileResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PlayerProfileResponse;
    static equals(a: PlayerProfileResponse | PlainMessage<PlayerProfileResponse> | undefined, b: PlayerProfileResponse | PlainMessage<PlayerProfileResponse> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.Lobby
 */
export declare class Lobby extends Message<Lobby> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    /**
     * @generated from field: string lobby_code = 2;
     */
    lobbyCode: string;
    /**
     * @generated from field: string lobby_name = 3;
     */
    lobbyName: string;
    /**
     * @generated from field: string host_user_id = 4;
     */
    hostUserId: string;
    /**
     * @generated from field: repeated string participants = 5;
     */
    participants: string[];
    /**
     * @generated from field: optional string game_id = 6;
     */
    gameId?: string;
    constructor(data?: PartialMessage<Lobby>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.Lobby";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Lobby;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Lobby;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Lobby;
    static equals(a: Lobby | PlainMessage<Lobby> | undefined, b: Lobby | PlainMessage<Lobby> | undefined): boolean;
}
/**
 * Message for create a lobby
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.CreateLobbyRequest
 */
export declare class CreateLobbyRequest extends Message<CreateLobbyRequest> {
    /**
     * @generated from field: string lobby_name = 1;
     */
    lobbyName: string;
    constructor(data?: PartialMessage<CreateLobbyRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.CreateLobbyRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateLobbyRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateLobbyRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateLobbyRequest;
    static equals(a: CreateLobbyRequest | PlainMessage<CreateLobbyRequest> | undefined, b: CreateLobbyRequest | PlainMessage<CreateLobbyRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.CreateLobbyResponse
 */
export declare class CreateLobbyResponse extends Message<CreateLobbyResponse> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    constructor(data?: PartialMessage<CreateLobbyResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.CreateLobbyResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): CreateLobbyResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): CreateLobbyResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): CreateLobbyResponse;
    static equals(a: CreateLobbyResponse | PlainMessage<CreateLobbyResponse> | undefined, b: CreateLobbyResponse | PlainMessage<CreateLobbyResponse> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.GetLobbyRequest
 */
export declare class GetLobbyRequest extends Message<GetLobbyRequest> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    constructor(data?: PartialMessage<GetLobbyRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GetLobbyRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetLobbyRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetLobbyRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetLobbyRequest;
    static equals(a: GetLobbyRequest | PlainMessage<GetLobbyRequest> | undefined, b: GetLobbyRequest | PlainMessage<GetLobbyRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.GetLobbyReply
 */
export declare class GetLobbyReply extends Message<GetLobbyReply> {
    /**
     * @generated from field: com.sweetloveinyourheart.kittens.clients.Lobby lobby = 1;
     */
    lobby?: Lobby;
    constructor(data?: PartialMessage<GetLobbyReply>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GetLobbyReply";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetLobbyReply;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetLobbyReply;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetLobbyReply;
    static equals(a: GetLobbyReply | PlainMessage<GetLobbyReply> | undefined, b: GetLobbyReply | PlainMessage<GetLobbyReply> | undefined): boolean;
}
/**
 * Message for join a lobby
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.JoinLobbyRequest
 */
export declare class JoinLobbyRequest extends Message<JoinLobbyRequest> {
    /**
     * @generated from field: string lobby_code = 1;
     */
    lobbyCode: string;
    constructor(data?: PartialMessage<JoinLobbyRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.JoinLobbyRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): JoinLobbyRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): JoinLobbyRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): JoinLobbyRequest;
    static equals(a: JoinLobbyRequest | PlainMessage<JoinLobbyRequest> | undefined, b: JoinLobbyRequest | PlainMessage<JoinLobbyRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.JoinLobbyResponse
 */
export declare class JoinLobbyResponse extends Message<JoinLobbyResponse> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    constructor(data?: PartialMessage<JoinLobbyResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.JoinLobbyResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): JoinLobbyResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): JoinLobbyResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): JoinLobbyResponse;
    static equals(a: JoinLobbyResponse | PlainMessage<JoinLobbyResponse> | undefined, b: JoinLobbyResponse | PlainMessage<JoinLobbyResponse> | undefined): boolean;
}
/**
 * Message for leave a lobby
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.LeaveLobbyRequest
 */
export declare class LeaveLobbyRequest extends Message<LeaveLobbyRequest> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    constructor(data?: PartialMessage<LeaveLobbyRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.LeaveLobbyRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LeaveLobbyRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LeaveLobbyRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LeaveLobbyRequest;
    static equals(a: LeaveLobbyRequest | PlainMessage<LeaveLobbyRequest> | undefined, b: LeaveLobbyRequest | PlainMessage<LeaveLobbyRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.LeaveLobbyResponse
 */
export declare class LeaveLobbyResponse extends Message<LeaveLobbyResponse> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    constructor(data?: PartialMessage<LeaveLobbyResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.LeaveLobbyResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LeaveLobbyResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LeaveLobbyResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LeaveLobbyResponse;
    static equals(a: LeaveLobbyResponse | PlainMessage<LeaveLobbyResponse> | undefined, b: LeaveLobbyResponse | PlainMessage<LeaveLobbyResponse> | undefined): boolean;
}
/**
 * Message for start a game
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.StartGameRequest
 */
export declare class StartGameRequest extends Message<StartGameRequest> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    constructor(data?: PartialMessage<StartGameRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.StartGameRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StartGameRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StartGameRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StartGameRequest;
    static equals(a: StartGameRequest | PlainMessage<StartGameRequest> | undefined, b: StartGameRequest | PlainMessage<StartGameRequest> | undefined): boolean;
}
