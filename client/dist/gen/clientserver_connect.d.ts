import { CreateLobbyRequest, CreateLobbyResponse, CreateNewGuestUserRequest, CreateNewGuestUserResponse, GetLobbyReply, GetLobbyRequest, GuestLoginRequest, GuestLoginResponse, JoinLobbyRequest, JoinLobbyResponse, LeaveLobbyRequest, LeaveLobbyResponse, PlayerProfileRequest, PlayerProfileResponse, StartGameRequest } from "./clientserver_pb.js";
import { Empty, MethodKind } from "@bufbuild/protobuf";
/**
 * @generated from service com.sweetloveinyourheart.kittens.clients.ClientServer
 */
export declare const ClientServer: {
    readonly typeName: "com.sweetloveinyourheart.kittens.clients.ClientServer";
    readonly methods: {
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
            readonly O: typeof PlayerProfileResponse;
            readonly kind: MethodKind.Unary;
        };
        /**
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayerProfile
         */
        readonly getPlayerProfile: {
            readonly name: "GetPlayerProfile";
            readonly I: typeof PlayerProfileRequest;
            readonly O: typeof PlayerProfileResponse;
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
         * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StartGame
         */
        readonly startGame: {
            readonly name: "StartGame";
            readonly I: typeof StartGameRequest;
            readonly O: typeof Empty;
            readonly kind: MethodKind.Unary;
        };
    };
};
