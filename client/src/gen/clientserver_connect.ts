// @generated by protoc-gen-connect-es v1.6.1 with parameter "target=ts"
// @generated from file clientserver.proto (package com.sweetloveinyourheart.kittens.clients, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { CreateLobbyRequest, CreateLobbyResponse, CreateNewGuestUserRequest, CreateNewGuestUserResponse, GetLobbyReply, GetLobbyRequest, GuestLoginRequest, GuestLoginResponse, JoinLobbyRequest, JoinLobbyResponse, LeaveLobbyRequest, LeaveLobbyResponse, PlayerProfileRequest, PlayerProfileResponse, StartGameRequest } from "./clientserver_pb.js";
import { Empty, MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service com.sweetloveinyourheart.kittens.clients.ClientServer
 */
export const ClientServer = {
  typeName: "com.sweetloveinyourheart.kittens.clients.ClientServer",
  methods: {
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.CreateNewGuestUser
     */
    createNewGuestUser: {
      name: "CreateNewGuestUser",
      I: CreateNewGuestUserRequest,
      O: CreateNewGuestUserResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GuestLogin
     */
    guestLogin: {
      name: "GuestLogin",
      I: GuestLoginRequest,
      O: GuestLoginResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetUserProfile
     */
    getUserProfile: {
      name: "GetUserProfile",
      I: Empty,
      O: PlayerProfileResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayerProfile
     */
    getPlayerProfile: {
      name: "GetPlayerProfile",
      I: PlayerProfileRequest,
      O: PlayerProfileResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.CreateLobby
     */
    createLobby: {
      name: "CreateLobby",
      I: CreateLobbyRequest,
      O: CreateLobbyResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StreamLobby
     */
    streamLobby: {
      name: "StreamLobby",
      I: GetLobbyRequest,
      O: GetLobbyReply,
      kind: MethodKind.ServerStreaming,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.JoinLobby
     */
    joinLobby: {
      name: "JoinLobby",
      I: JoinLobbyRequest,
      O: JoinLobbyResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.LeaveLobby
     */
    leaveLobby: {
      name: "LeaveLobby",
      I: LeaveLobbyRequest,
      O: LeaveLobbyResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StartGame
     */
    startGame: {
      name: "StartGame",
      I: StartGameRequest,
      O: Empty,
      kind: MethodKind.Unary,
    },
  }
} as const;

