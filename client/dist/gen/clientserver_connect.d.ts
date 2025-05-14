import { Empty, MethodKind } from "@bufbuild/protobuf";
import { CreateLobbyRequest, CreateLobbyResponse, CreateNewGuestUserRequest, CreateNewGuestUserResponse, GetGameMetaDataRequest, GetGameMetaDataResponse, GetLobbyReply, GetLobbyRequest, GiveCardRequest, GuestLoginRequest, GuestLoginResponse, JoinLobbyRequest, JoinLobbyResponse, LeaveLobbyRequest, LeaveLobbyResponse, PeekCardsRequest, PeekCardsResponse, PlayCardsRequest, PlayersProfileRequest, PlayersProfileResponse, RetrieveCardsDataResponse, SelectAffectedPlayerRequest, StartMatchRequest, StealCardRequest, StreamGameReply, StreamGameRequest, UserProfileResponse } from "./clientserver_pb.js";
/**
 * @generated from service com.sweetloveinyourheart.kittens.clients.ClientServer
 */
export declare const ClientServer: {
    readonly typeName: "com.sweetloveinyourheart.kittens.clients.ClientServer";
    readonly methods: {
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.RetrieveCardsData
         */
        readonly retrieveCardsData: {
            readonly name: "RetrieveCardsData";
            readonly I: typeof Empty;
            readonly O: typeof RetrieveCardsDataResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.CreateNewGuestUser
         */
        readonly createNewGuestUser: {
            readonly name: "CreateNewGuestUser";
            readonly I: typeof CreateNewGuestUserRequest;
            readonly O: typeof CreateNewGuestUserResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GuestLogin
         */
        readonly guestLogin: {
            readonly name: "GuestLogin";
            readonly I: typeof GuestLoginRequest;
            readonly O: typeof GuestLoginResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetUserProfile
         */
        readonly getUserProfile: {
            readonly name: "GetUserProfile";
            readonly I: typeof Empty;
            readonly O: typeof UserProfileResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayersProfile
         */
        readonly getPlayersProfile: {
            readonly name: "GetPlayersProfile";
            readonly I: typeof PlayersProfileRequest;
            readonly O: typeof PlayersProfileResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.CreateLobby
         */
        readonly createLobby: {
            readonly name: "CreateLobby";
            readonly I: typeof CreateLobbyRequest;
            readonly O: typeof CreateLobbyResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetLobby
         */
        readonly getLobby: {
            readonly name: "GetLobby";
            readonly I: typeof GetLobbyRequest;
            readonly O: typeof GetLobbyReply;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StreamLobby
         */
        readonly streamLobby: {
            readonly name: "StreamLobby";
            readonly I: typeof GetLobbyRequest;
            readonly O: typeof GetLobbyReply;
            readonly kind: MethodKind.ServerStreaming;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.JoinLobby
         */
        readonly joinLobby: {
            readonly name: "JoinLobby";
            readonly I: typeof JoinLobbyRequest;
            readonly O: typeof JoinLobbyResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.LeaveLobby
         */
        readonly leaveLobby: {
            readonly name: "LeaveLobby";
            readonly I: typeof LeaveLobbyRequest;
            readonly O: typeof LeaveLobbyResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StartMatch
         */
        readonly startMatch: {
            readonly name: "StartMatch";
            readonly I: typeof StartMatchRequest;
            readonly O: typeof Empty;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetGameMetaData
         */
        readonly getGameMetaData: {
            readonly name: "GetGameMetaData";
            readonly I: typeof GetGameMetaDataRequest;
            readonly O: typeof GetGameMetaDataResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StreamGame
         */
        readonly streamGame: {
            readonly name: "StreamGame";
            readonly I: typeof StreamGameRequest;
            readonly O: typeof StreamGameReply;
            readonly kind: MethodKind.ServerStreaming;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.PlayCards
         */
        readonly playCards: {
            readonly name: "PlayCards";
            readonly I: typeof PlayCardsRequest;
            readonly O: typeof Empty;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.PeekCards
         */
        readonly peekCards: {
            readonly name: "PeekCards";
            readonly I: typeof PeekCardsRequest;
            readonly O: typeof PeekCardsResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.SelectAffectedPlayer
         */
        readonly selectAffectedPlayer: {
            readonly name: "SelectAffectedPlayer";
            readonly I: typeof SelectAffectedPlayerRequest;
            readonly O: typeof Empty;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StealCard
         */
        readonly stealCard: {
            readonly name: "StealCard";
            readonly I: typeof StealCardRequest;
            readonly O: typeof Empty;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GiveCard
         */
        readonly giveCard: {
            readonly name: "GiveCard";
            readonly I: typeof GiveCardRequest;
            readonly O: typeof Empty;
            readonly kind: MethodKind.Unary;
        };
    };
};
