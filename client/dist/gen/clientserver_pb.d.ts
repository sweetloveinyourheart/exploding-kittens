import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.Card
 */
export declare class Card extends Message<Card> {
    /**
     * @generated from field: string card_id = 1;
     */
    cardId: string;
    /**
     * @generated from field: string name = 2;
     */
    name: string;
    /**
     * @generated from field: string code = 3;
     */
    code: string;
    /**
     * @generated from field: string description = 4;
     */
    description: string;
    constructor(data?: PartialMessage<Card>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.Card";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Card;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Card;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Card;
    static equals(a: Card | PlainMessage<Card> | undefined, b: Card | PlainMessage<Card> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.RetrieveCardsDataResponse
 */
export declare class RetrieveCardsDataResponse extends Message<RetrieveCardsDataResponse> {
    /**
     * @generated from field: repeated com.sweetloveinyourheart.kittens.clients.Card cards = 1;
     */
    cards: Card[];
    constructor(data?: PartialMessage<RetrieveCardsDataResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.RetrieveCardsDataResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): RetrieveCardsDataResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): RetrieveCardsDataResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): RetrieveCardsDataResponse;
    static equals(a: RetrieveCardsDataResponse | PlainMessage<RetrieveCardsDataResponse> | undefined, b: RetrieveCardsDataResponse | PlainMessage<RetrieveCardsDataResponse> | undefined): boolean;
}
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
 * @generated from message com.sweetloveinyourheart.kittens.clients.UserProfileResponse
 */
export declare class UserProfileResponse extends Message<UserProfileResponse> {
    /**
     * @generated from field: com.sweetloveinyourheart.kittens.clients.User user = 1;
     */
    user?: User;
    constructor(data?: PartialMessage<UserProfileResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.UserProfileResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UserProfileResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UserProfileResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UserProfileResponse;
    static equals(a: UserProfileResponse | PlainMessage<UserProfileResponse> | undefined, b: UserProfileResponse | PlainMessage<UserProfileResponse> | undefined): boolean;
}
/**
 * Message for players profile
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.PlayersProfileRequest
 */
export declare class PlayersProfileRequest extends Message<PlayersProfileRequest> {
    /**
     * Required: UUID of the guest user
     *
     * @generated from field: repeated string user_ids = 1;
     */
    userIds: string[];
    constructor(data?: PartialMessage<PlayersProfileRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.PlayersProfileRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PlayersProfileRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PlayersProfileRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PlayersProfileRequest;
    static equals(a: PlayersProfileRequest | PlainMessage<PlayersProfileRequest> | undefined, b: PlayersProfileRequest | PlainMessage<PlayersProfileRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.PlayersProfileResponse
 */
export declare class PlayersProfileResponse extends Message<PlayersProfileResponse> {
    /**
     * @generated from field: repeated com.sweetloveinyourheart.kittens.clients.User users = 1;
     */
    users: User[];
    constructor(data?: PartialMessage<PlayersProfileResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.PlayersProfileResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PlayersProfileResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PlayersProfileResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PlayersProfileResponse;
    static equals(a: PlayersProfileResponse | PlainMessage<PlayersProfileResponse> | undefined, b: PlayersProfileResponse | PlainMessage<PlayersProfileResponse> | undefined): boolean;
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
     * @generated from field: optional string match_id = 6;
     */
    matchId?: string;
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
 * Message for start a match
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.StartMatchRequest
 */
export declare class StartMatchRequest extends Message<StartMatchRequest> {
    /**
     * @generated from field: string lobby_id = 1;
     */
    lobbyId: string;
    constructor(data?: PartialMessage<StartMatchRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.StartMatchRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StartMatchRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StartMatchRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StartMatchRequest;
    static equals(a: StartMatchRequest | PlainMessage<StartMatchRequest> | undefined, b: StartMatchRequest | PlainMessage<StartMatchRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.Game
 */
export declare class Game extends Message<Game> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    /**
     * @generated from field: com.sweetloveinyourheart.kittens.clients.Game.Phase game_phase = 2;
     */
    gamePhase: Game_Phase;
    /**
     * @generated from field: string player_turn = 3;
     */
    playerTurn: string;
    /**
     * @generated from field: repeated com.sweetloveinyourheart.kittens.clients.Game.Player players = 4;
     */
    players: Game_Player[];
    /**
     * @generated from field: map<string, com.sweetloveinyourheart.kittens.clients.Game.PlayerHand> player_hands = 5;
     */
    playerHands: {
        [key: string]: Game_PlayerHand;
    };
    /**
     * @generated from field: com.sweetloveinyourheart.kittens.clients.Game.Desk desk = 6;
     */
    desk?: Game_Desk;
    /**
     * @generated from field: string executing_action = 7;
     */
    executingAction: string;
    /**
     * @generated from field: string affected_player = 8;
     */
    affectedPlayer: string;
    constructor(data?: PartialMessage<Game>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.Game";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Game;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Game;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Game;
    static equals(a: Game | PlainMessage<Game> | undefined, b: Game | PlainMessage<Game> | undefined): boolean;
}
/**
 * @generated from enum com.sweetloveinyourheart.kittens.clients.Game.Phase
 */
export declare enum Game_Phase {
    /**
     * Setting up players, shuffling and dealing cards, inserting Exploding Kittens and Defuse cards into the deck
     *
     * @generated from enum value: INITIALIZING = 0;
     */
    INITIALIZING = 0,
    /**
     * Active player begins their turn
     *
     * @generated from enum value: TURN_START = 1;
     */
    TURN_START = 1,
    /**
     * Player can play as many action cards as they want
     *
     * @generated from enum value: ACTION_PHASE = 2;
     */
    ACTION_PHASE = 2,
    /**
     * Player draws one card from the deck (mandatory if they didn't Skip/Attack)
     *
     * @generated from enum value: CARD_DRAWING = 3;
     */
    CARD_DRAWING = 3,
    /**
     * Finalize the turn, next player becomes active
     *
     * @generated from enum value: TURN_END = 4;
     */
    TURN_END = 4,
    /**
     * When only one player remains
     *
     * @generated from enum value: GAME_OVER = 5;
     */
    GAME_OVER = 5
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.Game.Player
 */
export declare class Game_Player extends Message<Game_Player> {
    /**
     * @generated from field: string player_id = 1;
     */
    playerId: string;
    /**
     * @generated from field: bool active = 2;
     */
    active: boolean;
    constructor(data?: PartialMessage<Game_Player>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.Game.Player";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Game_Player;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Game_Player;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Game_Player;
    static equals(a: Game_Player | PlainMessage<Game_Player> | undefined, b: Game_Player | PlainMessage<Game_Player> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.Game.PlayerHand
 */
export declare class Game_PlayerHand extends Message<Game_PlayerHand> {
    /**
     * @generated from field: int32 remaining_cards = 1;
     */
    remainingCards: number;
    /**
     * @generated from field: repeated string hands = 2;
     */
    hands: string[];
    constructor(data?: PartialMessage<Game_PlayerHand>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.Game.PlayerHand";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Game_PlayerHand;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Game_PlayerHand;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Game_PlayerHand;
    static equals(a: Game_PlayerHand | PlainMessage<Game_PlayerHand> | undefined, b: Game_PlayerHand | PlainMessage<Game_PlayerHand> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.Game.Desk
 */
export declare class Game_Desk extends Message<Game_Desk> {
    /**
     * @generated from field: string desk_id = 1;
     */
    deskId: string;
    /**
     * @generated from field: int32 remaining_cards = 2;
     */
    remainingCards: number;
    /**
     * @generated from field: repeated string discard_pile = 3;
     */
    discardPile: string[];
    constructor(data?: PartialMessage<Game_Desk>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.Game.Desk";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Game_Desk;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Game_Desk;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Game_Desk;
    static equals(a: Game_Desk | PlainMessage<Game_Desk> | undefined, b: Game_Desk | PlainMessage<Game_Desk> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.GameMetaData
 */
export declare class GameMetaData extends Message<GameMetaData> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    /**
     * @generated from field: repeated string players = 2;
     */
    players: string[];
    constructor(data?: PartialMessage<GameMetaData>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GameMetaData";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GameMetaData;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GameMetaData;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GameMetaData;
    static equals(a: GameMetaData | PlainMessage<GameMetaData> | undefined, b: GameMetaData | PlainMessage<GameMetaData> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.StreamGameRequest
 */
export declare class StreamGameRequest extends Message<StreamGameRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    constructor(data?: PartialMessage<StreamGameRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.StreamGameRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StreamGameRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StreamGameRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StreamGameRequest;
    static equals(a: StreamGameRequest | PlainMessage<StreamGameRequest> | undefined, b: StreamGameRequest | PlainMessage<StreamGameRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.StreamGameReply
 */
export declare class StreamGameReply extends Message<StreamGameReply> {
    /**
     * @generated from field: com.sweetloveinyourheart.kittens.clients.Game game_state = 1;
     */
    gameState?: Game;
    constructor(data?: PartialMessage<StreamGameReply>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.StreamGameReply";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StreamGameReply;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StreamGameReply;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StreamGameReply;
    static equals(a: StreamGameReply | PlainMessage<StreamGameReply> | undefined, b: StreamGameReply | PlainMessage<StreamGameReply> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.GetGameMetaDataRequest
 */
export declare class GetGameMetaDataRequest extends Message<GetGameMetaDataRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    constructor(data?: PartialMessage<GetGameMetaDataRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GetGameMetaDataRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetGameMetaDataRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetGameMetaDataRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetGameMetaDataRequest;
    static equals(a: GetGameMetaDataRequest | PlainMessage<GetGameMetaDataRequest> | undefined, b: GetGameMetaDataRequest | PlainMessage<GetGameMetaDataRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.GetGameMetaDataResponse
 */
export declare class GetGameMetaDataResponse extends Message<GetGameMetaDataResponse> {
    /**
     * @generated from field: com.sweetloveinyourheart.kittens.clients.GameMetaData meta = 1;
     */
    meta?: GameMetaData;
    constructor(data?: PartialMessage<GetGameMetaDataResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GetGameMetaDataResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetGameMetaDataResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetGameMetaDataResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetGameMetaDataResponse;
    static equals(a: GetGameMetaDataResponse | PlainMessage<GetGameMetaDataResponse> | undefined, b: GetGameMetaDataResponse | PlainMessage<GetGameMetaDataResponse> | undefined): boolean;
}
/**
 * Message for playing cards
 * This message is used to play cards in the game
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.PlayCardsRequest
 */
export declare class PlayCardsRequest extends Message<PlayCardsRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    /**
     * @generated from field: repeated string card_ids = 2;
     */
    cardIds: string[];
    constructor(data?: PartialMessage<PlayCardsRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.PlayCardsRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PlayCardsRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PlayCardsRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PlayCardsRequest;
    static equals(a: PlayCardsRequest | PlainMessage<PlayCardsRequest> | undefined, b: PlayCardsRequest | PlainMessage<PlayCardsRequest> | undefined): boolean;
}
/**
 * Message for peeking cards
 * This message is used to peek at the top card of the deck
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.PeekCardsRequest
 */
export declare class PeekCardsRequest extends Message<PeekCardsRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    /**
     * @generated from field: string desk_id = 2;
     */
    deskId: string;
    constructor(data?: PartialMessage<PeekCardsRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.PeekCardsRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PeekCardsRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PeekCardsRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PeekCardsRequest;
    static equals(a: PeekCardsRequest | PlainMessage<PeekCardsRequest> | undefined, b: PeekCardsRequest | PlainMessage<PeekCardsRequest> | undefined): boolean;
}
/**
 * @generated from message com.sweetloveinyourheart.kittens.clients.PeekCardsResponse
 */
export declare class PeekCardsResponse extends Message<PeekCardsResponse> {
    /**
     * @generated from field: repeated string card_ids = 1;
     */
    cardIds: string[];
    constructor(data?: PartialMessage<PeekCardsResponse>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.PeekCardsResponse";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PeekCardsResponse;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PeekCardsResponse;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PeekCardsResponse;
    static equals(a: PeekCardsResponse | PlainMessage<PeekCardsResponse> | undefined, b: PeekCardsResponse | PlainMessage<PeekCardsResponse> | undefined): boolean;
}
/**
 * Message for selecting affected players
 * This message is used to select affected players in the game
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.SelectAffectedPlayerRequest
 */
export declare class SelectAffectedPlayerRequest extends Message<SelectAffectedPlayerRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    /**
     * @generated from field: string player_id = 2;
     */
    playerId: string;
    constructor(data?: PartialMessage<SelectAffectedPlayerRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.SelectAffectedPlayerRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SelectAffectedPlayerRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SelectAffectedPlayerRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SelectAffectedPlayerRequest;
    static equals(a: SelectAffectedPlayerRequest | PlainMessage<SelectAffectedPlayerRequest> | undefined, b: SelectAffectedPlayerRequest | PlainMessage<SelectAffectedPlayerRequest> | undefined): boolean;
}
/**
 * Message for stealing a card
 * This message is used to steal a card from another player
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.StealCardRequest
 */
export declare class StealCardRequest extends Message<StealCardRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    /**
     * @generated from field: optional string card_id = 2;
     */
    cardId?: string;
    /**
     * @generated from field: optional int32 card_index = 3;
     */
    cardIndex?: number;
    constructor(data?: PartialMessage<StealCardRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.StealCardRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StealCardRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StealCardRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StealCardRequest;
    static equals(a: StealCardRequest | PlainMessage<StealCardRequest> | undefined, b: StealCardRequest | PlainMessage<StealCardRequest> | undefined): boolean;
}
/**
 * Message for giving a card
 * This message is used to give a card to another player
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.GiveCardRequest
 */
export declare class GiveCardRequest extends Message<GiveCardRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    /**
     * @generated from field: string card_id = 2;
     */
    cardId: string;
    constructor(data?: PartialMessage<GiveCardRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.GiveCardRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GiveCardRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GiveCardRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GiveCardRequest;
    static equals(a: GiveCardRequest | PlainMessage<GiveCardRequest> | undefined, b: GiveCardRequest | PlainMessage<GiveCardRequest> | undefined): boolean;
}
/**
 * Message for drawing cards
 * This message is used to draw cards from the deck
 *
 * @generated from message com.sweetloveinyourheart.kittens.clients.DrawCardRequest
 */
export declare class DrawCardRequest extends Message<DrawCardRequest> {
    /**
     * @generated from field: string game_id = 1;
     */
    gameId: string;
    constructor(data?: PartialMessage<DrawCardRequest>);
    static readonly runtime: typeof proto3;
    static readonly typeName = "com.sweetloveinyourheart.kittens.clients.DrawCardRequest";
    static readonly fields: FieldList;
    static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DrawCardRequest;
    static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DrawCardRequest;
    static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DrawCardRequest;
    static equals(a: DrawCardRequest | PlainMessage<DrawCardRequest> | undefined, b: DrawCardRequest | PlainMessage<DrawCardRequest> | undefined): boolean;
}
