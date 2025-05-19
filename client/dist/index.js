// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/assert.js
function assert(condition, msg) {
  if (!condition) {
    throw new Error(msg);
  }
}
var FLOAT32_MAX = 34028234663852886e22;
var FLOAT32_MIN = -34028234663852886e22;
var UINT32_MAX = 4294967295;
var INT32_MAX = 2147483647;
var INT32_MIN = -2147483648;
function assertInt32(arg) {
  if (typeof arg !== "number")
    throw new Error("invalid int 32: " + typeof arg);
  if (!Number.isInteger(arg) || arg > INT32_MAX || arg < INT32_MIN)
    throw new Error("invalid int 32: " + arg);
}
function assertUInt32(arg) {
  if (typeof arg !== "number")
    throw new Error("invalid uint 32: " + typeof arg);
  if (!Number.isInteger(arg) || arg > UINT32_MAX || arg < 0)
    throw new Error("invalid uint 32: " + arg);
}
function assertFloat32(arg) {
  if (typeof arg !== "number")
    throw new Error("invalid float 32: " + typeof arg);
  if (!Number.isFinite(arg))
    return;
  if (arg > FLOAT32_MAX || arg < FLOAT32_MIN)
    throw new Error("invalid float 32: " + arg);
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/enum.js
var enumTypeSymbol = Symbol("@bufbuild/protobuf/enum-type");
function getEnumType(enumObject) {
  const t = enumObject[enumTypeSymbol];
  assert(t, "missing enum type on enum object");
  return t;
}
function setEnumType(enumObject, typeName, values, opt) {
  enumObject[enumTypeSymbol] = makeEnumType(typeName, values.map((v) => ({
    no: v.no,
    name: v.name,
    localName: enumObject[v.no]
  })), opt);
}
function makeEnumType(typeName, values, _opt) {
  const names = /* @__PURE__ */ Object.create(null);
  const numbers = /* @__PURE__ */ Object.create(null);
  const normalValues = [];
  for (const value of values) {
    const n = normalizeEnumValue(value);
    normalValues.push(n);
    names[value.name] = n;
    numbers[value.no] = n;
  }
  return {
    typeName,
    values: normalValues,
    // We do not surface options at this time
    // options: opt?.options ?? Object.create(null),
    findName(name) {
      return names[name];
    },
    findNumber(no) {
      return numbers[no];
    }
  };
}
function makeEnum(typeName, values, opt) {
  const enumObject = {};
  for (const value of values) {
    const n = normalizeEnumValue(value);
    enumObject[n.localName] = n.no;
    enumObject[n.no] = n.localName;
  }
  setEnumType(enumObject, typeName, values, opt);
  return enumObject;
}
function normalizeEnumValue(value) {
  if ("localName" in value) {
    return value;
  }
  return Object.assign(Object.assign({}, value), { localName: value.name });
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/message.js
var Message = class {
  /**
   * Compare with a message of the same type.
   * Note that this function disregards extensions and unknown fields.
   */
  equals(other) {
    return this.getType().runtime.util.equals(this.getType(), this, other);
  }
  /**
   * Create a deep copy.
   */
  clone() {
    return this.getType().runtime.util.clone(this);
  }
  /**
   * Parse from binary data, merging fields.
   *
   * Repeated fields are appended. Map entries are added, overwriting
   * existing keys.
   *
   * If a message field is already present, it will be merged with the
   * new data.
   */
  fromBinary(bytes, options) {
    const type = this.getType(), format = type.runtime.bin, opt = format.makeReadOptions(options);
    format.readMessage(this, opt.readerFactory(bytes), bytes.byteLength, opt);
    return this;
  }
  /**
   * Parse a message from a JSON value.
   */
  fromJson(jsonValue, options) {
    const type = this.getType(), format = type.runtime.json, opt = format.makeReadOptions(options);
    format.readMessage(type, jsonValue, opt, this);
    return this;
  }
  /**
   * Parse a message from a JSON string.
   */
  fromJsonString(jsonString, options) {
    let json;
    try {
      json = JSON.parse(jsonString);
    } catch (e) {
      throw new Error(`cannot decode ${this.getType().typeName} from JSON: ${e instanceof Error ? e.message : String(e)}`);
    }
    return this.fromJson(json, options);
  }
  /**
   * Serialize the message to binary data.
   */
  toBinary(options) {
    const type = this.getType(), bin = type.runtime.bin, opt = bin.makeWriteOptions(options), writer = opt.writerFactory();
    bin.writeMessage(this, writer, opt);
    return writer.finish();
  }
  /**
   * Serialize the message to a JSON value, a JavaScript value that can be
   * passed to JSON.stringify().
   */
  toJson(options) {
    const type = this.getType(), json = type.runtime.json, opt = json.makeWriteOptions(options);
    return json.writeMessage(this, opt);
  }
  /**
   * Serialize the message to a JSON string.
   */
  toJsonString(options) {
    var _a;
    const value = this.toJson(options);
    return JSON.stringify(value, null, (_a = options === null || options === void 0 ? void 0 : options.prettySpaces) !== null && _a !== void 0 ? _a : 0);
  }
  /**
   * Override for serialization behavior. This will be invoked when calling
   * JSON.stringify on this message (i.e. JSON.stringify(msg)).
   *
   * Note that this will not serialize google.protobuf.Any with a packed
   * message because the protobuf JSON format specifies that it needs to be
   * unpacked, and this is only possible with a type registry to look up the
   * message type.  As a result, attempting to serialize a message with this
   * type will throw an Error.
   *
   * This method is protected because you should not need to invoke it
   * directly -- instead use JSON.stringify or toJsonString for
   * stringified JSON.  Alternatively, if actual JSON is desired, you should
   * use toJson.
   */
  toJSON() {
    return this.toJson({
      emitDefaultValues: true
    });
  }
  /**
   * Retrieve the MessageType of this message - a singleton that represents
   * the protobuf message declaration and provides metadata for reflection-
   * based operations.
   */
  getType() {
    return Object.getPrototypeOf(this).constructor;
  }
};

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/message-type.js
function makeMessageType(runtime, typeName, fields, opt) {
  var _a;
  const localName = (_a = opt === null || opt === void 0 ? void 0 : opt.localName) !== null && _a !== void 0 ? _a : typeName.substring(typeName.lastIndexOf(".") + 1);
  const type = {
    [localName]: function(data) {
      runtime.util.initFields(this);
      runtime.util.initPartial(data, this);
    }
  }[localName];
  Object.setPrototypeOf(type.prototype, new Message());
  Object.assign(type, {
    runtime,
    typeName,
    fields: runtime.util.newFieldList(fields),
    fromBinary(bytes, options) {
      return new type().fromBinary(bytes, options);
    },
    fromJson(jsonValue, options) {
      return new type().fromJson(jsonValue, options);
    },
    fromJsonString(jsonString, options) {
      return new type().fromJsonString(jsonString, options);
    },
    equals(a, b) {
      return runtime.util.equals(type, a, b);
    }
  });
  return type;
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/google/varint.js
function varint64read() {
  let lowBits = 0;
  let highBits = 0;
  for (let shift = 0; shift < 28; shift += 7) {
    let b = this.buf[this.pos++];
    lowBits |= (b & 127) << shift;
    if ((b & 128) == 0) {
      this.assertBounds();
      return [lowBits, highBits];
    }
  }
  let middleByte = this.buf[this.pos++];
  lowBits |= (middleByte & 15) << 28;
  highBits = (middleByte & 112) >> 4;
  if ((middleByte & 128) == 0) {
    this.assertBounds();
    return [lowBits, highBits];
  }
  for (let shift = 3; shift <= 31; shift += 7) {
    let b = this.buf[this.pos++];
    highBits |= (b & 127) << shift;
    if ((b & 128) == 0) {
      this.assertBounds();
      return [lowBits, highBits];
    }
  }
  throw new Error("invalid varint");
}
function varint64write(lo, hi, bytes) {
  for (let i = 0; i < 28; i = i + 7) {
    const shift = lo >>> i;
    const hasNext = !(shift >>> 7 == 0 && hi == 0);
    const byte = (hasNext ? shift | 128 : shift) & 255;
    bytes.push(byte);
    if (!hasNext) {
      return;
    }
  }
  const splitBits = lo >>> 28 & 15 | (hi & 7) << 4;
  const hasMoreBits = !(hi >> 3 == 0);
  bytes.push((hasMoreBits ? splitBits | 128 : splitBits) & 255);
  if (!hasMoreBits) {
    return;
  }
  for (let i = 3; i < 31; i = i + 7) {
    const shift = hi >>> i;
    const hasNext = !(shift >>> 7 == 0);
    const byte = (hasNext ? shift | 128 : shift) & 255;
    bytes.push(byte);
    if (!hasNext) {
      return;
    }
  }
  bytes.push(hi >>> 31 & 1);
}
var TWO_PWR_32_DBL = 4294967296;
function int64FromString(dec) {
  const minus = dec[0] === "-";
  if (minus) {
    dec = dec.slice(1);
  }
  const base = 1e6;
  let lowBits = 0;
  let highBits = 0;
  function add1e6digit(begin, end) {
    const digit1e6 = Number(dec.slice(begin, end));
    highBits *= base;
    lowBits = lowBits * base + digit1e6;
    if (lowBits >= TWO_PWR_32_DBL) {
      highBits = highBits + (lowBits / TWO_PWR_32_DBL | 0);
      lowBits = lowBits % TWO_PWR_32_DBL;
    }
  }
  add1e6digit(-24, -18);
  add1e6digit(-18, -12);
  add1e6digit(-12, -6);
  add1e6digit(-6);
  return minus ? negate(lowBits, highBits) : newBits(lowBits, highBits);
}
function int64ToString(lo, hi) {
  let bits = newBits(lo, hi);
  const negative = bits.hi & 2147483648;
  if (negative) {
    bits = negate(bits.lo, bits.hi);
  }
  const result = uInt64ToString(bits.lo, bits.hi);
  return negative ? "-" + result : result;
}
function uInt64ToString(lo, hi) {
  ({ lo, hi } = toUnsigned(lo, hi));
  if (hi <= 2097151) {
    return String(TWO_PWR_32_DBL * hi + lo);
  }
  const low = lo & 16777215;
  const mid = (lo >>> 24 | hi << 8) & 16777215;
  const high = hi >> 16 & 65535;
  let digitA = low + mid * 6777216 + high * 6710656;
  let digitB = mid + high * 8147497;
  let digitC = high * 2;
  const base = 1e7;
  if (digitA >= base) {
    digitB += Math.floor(digitA / base);
    digitA %= base;
  }
  if (digitB >= base) {
    digitC += Math.floor(digitB / base);
    digitB %= base;
  }
  return digitC.toString() + decimalFrom1e7WithLeadingZeros(digitB) + decimalFrom1e7WithLeadingZeros(digitA);
}
function toUnsigned(lo, hi) {
  return { lo: lo >>> 0, hi: hi >>> 0 };
}
function newBits(lo, hi) {
  return { lo: lo | 0, hi: hi | 0 };
}
function negate(lowBits, highBits) {
  highBits = ~highBits;
  if (lowBits) {
    lowBits = ~lowBits + 1;
  } else {
    highBits += 1;
  }
  return newBits(lowBits, highBits);
}
var decimalFrom1e7WithLeadingZeros = (digit1e7) => {
  const partial = String(digit1e7);
  return "0000000".slice(partial.length) + partial;
};
function varint32write(value, bytes) {
  if (value >= 0) {
    while (value > 127) {
      bytes.push(value & 127 | 128);
      value = value >>> 7;
    }
    bytes.push(value);
  } else {
    for (let i = 0; i < 9; i++) {
      bytes.push(value & 127 | 128);
      value = value >> 7;
    }
    bytes.push(1);
  }
}
function varint32read() {
  let b = this.buf[this.pos++];
  let result = b & 127;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 127) << 7;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 127) << 14;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 127) << 21;
  if ((b & 128) == 0) {
    this.assertBounds();
    return result;
  }
  b = this.buf[this.pos++];
  result |= (b & 15) << 28;
  for (let readBytes = 5; (b & 128) !== 0 && readBytes < 10; readBytes++)
    b = this.buf[this.pos++];
  if ((b & 128) != 0)
    throw new Error("invalid varint");
  this.assertBounds();
  return result >>> 0;
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/proto-int64.js
function makeInt64Support() {
  const dv = new DataView(new ArrayBuffer(8));
  const ok = typeof BigInt === "function" && typeof dv.getBigInt64 === "function" && typeof dv.getBigUint64 === "function" && typeof dv.setBigInt64 === "function" && typeof dv.setBigUint64 === "function" && (typeof process != "object" || typeof process.env != "object" || process.env.BUF_BIGINT_DISABLE !== "1");
  if (ok) {
    const MIN = BigInt("-9223372036854775808"), MAX = BigInt("9223372036854775807"), UMIN = BigInt("0"), UMAX = BigInt("18446744073709551615");
    return {
      zero: BigInt(0),
      supported: true,
      parse(value) {
        const bi = typeof value == "bigint" ? value : BigInt(value);
        if (bi > MAX || bi < MIN) {
          throw new Error(`int64 invalid: ${value}`);
        }
        return bi;
      },
      uParse(value) {
        const bi = typeof value == "bigint" ? value : BigInt(value);
        if (bi > UMAX || bi < UMIN) {
          throw new Error(`uint64 invalid: ${value}`);
        }
        return bi;
      },
      enc(value) {
        dv.setBigInt64(0, this.parse(value), true);
        return {
          lo: dv.getInt32(0, true),
          hi: dv.getInt32(4, true)
        };
      },
      uEnc(value) {
        dv.setBigInt64(0, this.uParse(value), true);
        return {
          lo: dv.getInt32(0, true),
          hi: dv.getInt32(4, true)
        };
      },
      dec(lo, hi) {
        dv.setInt32(0, lo, true);
        dv.setInt32(4, hi, true);
        return dv.getBigInt64(0, true);
      },
      uDec(lo, hi) {
        dv.setInt32(0, lo, true);
        dv.setInt32(4, hi, true);
        return dv.getBigUint64(0, true);
      }
    };
  }
  const assertInt64String = (value) => assert(/^-?[0-9]+$/.test(value), `int64 invalid: ${value}`);
  const assertUInt64String = (value) => assert(/^[0-9]+$/.test(value), `uint64 invalid: ${value}`);
  return {
    zero: "0",
    supported: false,
    parse(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertInt64String(value);
      return value;
    },
    uParse(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertUInt64String(value);
      return value;
    },
    enc(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertInt64String(value);
      return int64FromString(value);
    },
    uEnc(value) {
      if (typeof value != "string") {
        value = value.toString();
      }
      assertUInt64String(value);
      return int64FromString(value);
    },
    dec(lo, hi) {
      return int64ToString(lo, hi);
    },
    uDec(lo, hi) {
      return uInt64ToString(lo, hi);
    }
  };
}
var protoInt64 = makeInt64Support();

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/scalar.js
var ScalarType;
(function(ScalarType2) {
  ScalarType2[ScalarType2["DOUBLE"] = 1] = "DOUBLE";
  ScalarType2[ScalarType2["FLOAT"] = 2] = "FLOAT";
  ScalarType2[ScalarType2["INT64"] = 3] = "INT64";
  ScalarType2[ScalarType2["UINT64"] = 4] = "UINT64";
  ScalarType2[ScalarType2["INT32"] = 5] = "INT32";
  ScalarType2[ScalarType2["FIXED64"] = 6] = "FIXED64";
  ScalarType2[ScalarType2["FIXED32"] = 7] = "FIXED32";
  ScalarType2[ScalarType2["BOOL"] = 8] = "BOOL";
  ScalarType2[ScalarType2["STRING"] = 9] = "STRING";
  ScalarType2[ScalarType2["BYTES"] = 12] = "BYTES";
  ScalarType2[ScalarType2["UINT32"] = 13] = "UINT32";
  ScalarType2[ScalarType2["SFIXED32"] = 15] = "SFIXED32";
  ScalarType2[ScalarType2["SFIXED64"] = 16] = "SFIXED64";
  ScalarType2[ScalarType2["SINT32"] = 17] = "SINT32";
  ScalarType2[ScalarType2["SINT64"] = 18] = "SINT64";
})(ScalarType || (ScalarType = {}));
var LongType;
(function(LongType2) {
  LongType2[LongType2["BIGINT"] = 0] = "BIGINT";
  LongType2[LongType2["STRING"] = 1] = "STRING";
})(LongType || (LongType = {}));

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/scalars.js
function scalarEquals(type, a, b) {
  if (a === b) {
    return true;
  }
  if (type == ScalarType.BYTES) {
    if (!(a instanceof Uint8Array) || !(b instanceof Uint8Array)) {
      return false;
    }
    if (a.length !== b.length) {
      return false;
    }
    for (let i = 0; i < a.length; i++) {
      if (a[i] !== b[i]) {
        return false;
      }
    }
    return true;
  }
  switch (type) {
    case ScalarType.UINT64:
    case ScalarType.FIXED64:
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      return a == b;
  }
  return false;
}
function scalarZeroValue(type, longType) {
  switch (type) {
    case ScalarType.BOOL:
      return false;
    case ScalarType.UINT64:
    case ScalarType.FIXED64:
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      return longType == 0 ? protoInt64.zero : "0";
    case ScalarType.DOUBLE:
    case ScalarType.FLOAT:
      return 0;
    case ScalarType.BYTES:
      return new Uint8Array(0);
    case ScalarType.STRING:
      return "";
    default:
      return 0;
  }
}
function isScalarZeroValue(type, value) {
  switch (type) {
    case ScalarType.BOOL:
      return value === false;
    case ScalarType.STRING:
      return value === "";
    case ScalarType.BYTES:
      return value instanceof Uint8Array && !value.byteLength;
    default:
      return value == 0;
  }
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/binary-encoding.js
var WireType;
(function(WireType2) {
  WireType2[WireType2["Varint"] = 0] = "Varint";
  WireType2[WireType2["Bit64"] = 1] = "Bit64";
  WireType2[WireType2["LengthDelimited"] = 2] = "LengthDelimited";
  WireType2[WireType2["StartGroup"] = 3] = "StartGroup";
  WireType2[WireType2["EndGroup"] = 4] = "EndGroup";
  WireType2[WireType2["Bit32"] = 5] = "Bit32";
})(WireType || (WireType = {}));
var BinaryWriter = class {
  constructor(textEncoder) {
    this.stack = [];
    this.textEncoder = textEncoder !== null && textEncoder !== void 0 ? textEncoder : new TextEncoder();
    this.chunks = [];
    this.buf = [];
  }
  /**
   * Return all bytes written and reset this writer.
   */
  finish() {
    this.chunks.push(new Uint8Array(this.buf));
    let len = 0;
    for (let i = 0; i < this.chunks.length; i++)
      len += this.chunks[i].length;
    let bytes = new Uint8Array(len);
    let offset = 0;
    for (let i = 0; i < this.chunks.length; i++) {
      bytes.set(this.chunks[i], offset);
      offset += this.chunks[i].length;
    }
    this.chunks = [];
    return bytes;
  }
  /**
   * Start a new fork for length-delimited data like a message
   * or a packed repeated field.
   *
   * Must be joined later with `join()`.
   */
  fork() {
    this.stack.push({ chunks: this.chunks, buf: this.buf });
    this.chunks = [];
    this.buf = [];
    return this;
  }
  /**
   * Join the last fork. Write its length and bytes, then
   * return to the previous state.
   */
  join() {
    let chunk = this.finish();
    let prev = this.stack.pop();
    if (!prev)
      throw new Error("invalid state, fork stack empty");
    this.chunks = prev.chunks;
    this.buf = prev.buf;
    this.uint32(chunk.byteLength);
    return this.raw(chunk);
  }
  /**
   * Writes a tag (field number and wire type).
   *
   * Equivalent to `uint32( (fieldNo << 3 | type) >>> 0 )`.
   *
   * Generated code should compute the tag ahead of time and call `uint32()`.
   */
  tag(fieldNo, type) {
    return this.uint32((fieldNo << 3 | type) >>> 0);
  }
  /**
   * Write a chunk of raw bytes.
   */
  raw(chunk) {
    if (this.buf.length) {
      this.chunks.push(new Uint8Array(this.buf));
      this.buf = [];
    }
    this.chunks.push(chunk);
    return this;
  }
  /**
   * Write a `uint32` value, an unsigned 32 bit varint.
   */
  uint32(value) {
    assertUInt32(value);
    while (value > 127) {
      this.buf.push(value & 127 | 128);
      value = value >>> 7;
    }
    this.buf.push(value);
    return this;
  }
  /**
   * Write a `int32` value, a signed 32 bit varint.
   */
  int32(value) {
    assertInt32(value);
    varint32write(value, this.buf);
    return this;
  }
  /**
   * Write a `bool` value, a variant.
   */
  bool(value) {
    this.buf.push(value ? 1 : 0);
    return this;
  }
  /**
   * Write a `bytes` value, length-delimited arbitrary data.
   */
  bytes(value) {
    this.uint32(value.byteLength);
    return this.raw(value);
  }
  /**
   * Write a `string` value, length-delimited data converted to UTF-8 text.
   */
  string(value) {
    let chunk = this.textEncoder.encode(value);
    this.uint32(chunk.byteLength);
    return this.raw(chunk);
  }
  /**
   * Write a `float` value, 32-bit floating point number.
   */
  float(value) {
    assertFloat32(value);
    let chunk = new Uint8Array(4);
    new DataView(chunk.buffer).setFloat32(0, value, true);
    return this.raw(chunk);
  }
  /**
   * Write a `double` value, a 64-bit floating point number.
   */
  double(value) {
    let chunk = new Uint8Array(8);
    new DataView(chunk.buffer).setFloat64(0, value, true);
    return this.raw(chunk);
  }
  /**
   * Write a `fixed32` value, an unsigned, fixed-length 32-bit integer.
   */
  fixed32(value) {
    assertUInt32(value);
    let chunk = new Uint8Array(4);
    new DataView(chunk.buffer).setUint32(0, value, true);
    return this.raw(chunk);
  }
  /**
   * Write a `sfixed32` value, a signed, fixed-length 32-bit integer.
   */
  sfixed32(value) {
    assertInt32(value);
    let chunk = new Uint8Array(4);
    new DataView(chunk.buffer).setInt32(0, value, true);
    return this.raw(chunk);
  }
  /**
   * Write a `sint32` value, a signed, zigzag-encoded 32-bit varint.
   */
  sint32(value) {
    assertInt32(value);
    value = (value << 1 ^ value >> 31) >>> 0;
    varint32write(value, this.buf);
    return this;
  }
  /**
   * Write a `fixed64` value, a signed, fixed-length 64-bit integer.
   */
  sfixed64(value) {
    let chunk = new Uint8Array(8), view = new DataView(chunk.buffer), tc = protoInt64.enc(value);
    view.setInt32(0, tc.lo, true);
    view.setInt32(4, tc.hi, true);
    return this.raw(chunk);
  }
  /**
   * Write a `fixed64` value, an unsigned, fixed-length 64 bit integer.
   */
  fixed64(value) {
    let chunk = new Uint8Array(8), view = new DataView(chunk.buffer), tc = protoInt64.uEnc(value);
    view.setInt32(0, tc.lo, true);
    view.setInt32(4, tc.hi, true);
    return this.raw(chunk);
  }
  /**
   * Write a `int64` value, a signed 64-bit varint.
   */
  int64(value) {
    let tc = protoInt64.enc(value);
    varint64write(tc.lo, tc.hi, this.buf);
    return this;
  }
  /**
   * Write a `sint64` value, a signed, zig-zag-encoded 64-bit varint.
   */
  sint64(value) {
    let tc = protoInt64.enc(value), sign = tc.hi >> 31, lo = tc.lo << 1 ^ sign, hi = (tc.hi << 1 | tc.lo >>> 31) ^ sign;
    varint64write(lo, hi, this.buf);
    return this;
  }
  /**
   * Write a `uint64` value, an unsigned 64-bit varint.
   */
  uint64(value) {
    let tc = protoInt64.uEnc(value);
    varint64write(tc.lo, tc.hi, this.buf);
    return this;
  }
};
var BinaryReader = class {
  constructor(buf, textDecoder) {
    this.varint64 = varint64read;
    this.uint32 = varint32read;
    this.buf = buf;
    this.len = buf.length;
    this.pos = 0;
    this.view = new DataView(buf.buffer, buf.byteOffset, buf.byteLength);
    this.textDecoder = textDecoder !== null && textDecoder !== void 0 ? textDecoder : new TextDecoder();
  }
  /**
   * Reads a tag - field number and wire type.
   */
  tag() {
    let tag = this.uint32(), fieldNo = tag >>> 3, wireType = tag & 7;
    if (fieldNo <= 0 || wireType < 0 || wireType > 5)
      throw new Error("illegal tag: field no " + fieldNo + " wire type " + wireType);
    return [fieldNo, wireType];
  }
  /**
   * Skip one element and return the skipped data.
   *
   * When skipping StartGroup, provide the tags field number to check for
   * matching field number in the EndGroup tag.
   */
  skip(wireType, fieldNo) {
    let start = this.pos;
    switch (wireType) {
      case WireType.Varint:
        while (this.buf[this.pos++] & 128) {
        }
        break;
      // eslint-disable-next-line
      // @ts-ignore TS7029: Fallthrough case in switch
      case WireType.Bit64:
        this.pos += 4;
      // eslint-disable-next-line
      // @ts-ignore TS7029: Fallthrough case in switch
      case WireType.Bit32:
        this.pos += 4;
        break;
      case WireType.LengthDelimited:
        let len = this.uint32();
        this.pos += len;
        break;
      case WireType.StartGroup:
        for (; ; ) {
          const [fn, wt] = this.tag();
          if (wt === WireType.EndGroup) {
            if (fieldNo !== void 0 && fn !== fieldNo) {
              throw new Error("invalid end group tag");
            }
            break;
          }
          this.skip(wt, fn);
        }
        break;
      default:
        throw new Error("cant skip wire type " + wireType);
    }
    this.assertBounds();
    return this.buf.subarray(start, this.pos);
  }
  /**
   * Throws error if position in byte array is out of range.
   */
  assertBounds() {
    if (this.pos > this.len)
      throw new RangeError("premature EOF");
  }
  /**
   * Read a `int32` field, a signed 32 bit varint.
   */
  int32() {
    return this.uint32() | 0;
  }
  /**
   * Read a `sint32` field, a signed, zigzag-encoded 32-bit varint.
   */
  sint32() {
    let zze = this.uint32();
    return zze >>> 1 ^ -(zze & 1);
  }
  /**
   * Read a `int64` field, a signed 64-bit varint.
   */
  int64() {
    return protoInt64.dec(...this.varint64());
  }
  /**
   * Read a `uint64` field, an unsigned 64-bit varint.
   */
  uint64() {
    return protoInt64.uDec(...this.varint64());
  }
  /**
   * Read a `sint64` field, a signed, zig-zag-encoded 64-bit varint.
   */
  sint64() {
    let [lo, hi] = this.varint64();
    let s = -(lo & 1);
    lo = (lo >>> 1 | (hi & 1) << 31) ^ s;
    hi = hi >>> 1 ^ s;
    return protoInt64.dec(lo, hi);
  }
  /**
   * Read a `bool` field, a variant.
   */
  bool() {
    let [lo, hi] = this.varint64();
    return lo !== 0 || hi !== 0;
  }
  /**
   * Read a `fixed32` field, an unsigned, fixed-length 32-bit integer.
   */
  fixed32() {
    return this.view.getUint32((this.pos += 4) - 4, true);
  }
  /**
   * Read a `sfixed32` field, a signed, fixed-length 32-bit integer.
   */
  sfixed32() {
    return this.view.getInt32((this.pos += 4) - 4, true);
  }
  /**
   * Read a `fixed64` field, an unsigned, fixed-length 64 bit integer.
   */
  fixed64() {
    return protoInt64.uDec(this.sfixed32(), this.sfixed32());
  }
  /**
   * Read a `fixed64` field, a signed, fixed-length 64-bit integer.
   */
  sfixed64() {
    return protoInt64.dec(this.sfixed32(), this.sfixed32());
  }
  /**
   * Read a `float` field, 32-bit floating point number.
   */
  float() {
    return this.view.getFloat32((this.pos += 4) - 4, true);
  }
  /**
   * Read a `double` field, a 64-bit floating point number.
   */
  double() {
    return this.view.getFloat64((this.pos += 8) - 8, true);
  }
  /**
   * Read a `bytes` field, length-delimited arbitrary data.
   */
  bytes() {
    let len = this.uint32(), start = this.pos;
    this.pos += len;
    this.assertBounds();
    return this.buf.subarray(start, start + len);
  }
  /**
   * Read a `string` field, length-delimited data converted to UTF-8 text.
   */
  string() {
    return this.textDecoder.decode(this.bytes());
  }
};

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/extensions.js
function makeExtension(runtime, typeName, extendee, field) {
  let fi;
  return {
    typeName,
    extendee,
    get field() {
      if (!fi) {
        const i = typeof field == "function" ? field() : field;
        i.name = typeName.split(".").pop();
        i.jsonName = `[${typeName}]`;
        fi = runtime.util.newFieldList([i]).list()[0];
      }
      return fi;
    },
    runtime
  };
}
function createExtensionContainer(extension) {
  const localName = extension.field.localName;
  const container = /* @__PURE__ */ Object.create(null);
  container[localName] = initExtensionField(extension);
  return [container, () => container[localName]];
}
function initExtensionField(ext) {
  const field = ext.field;
  if (field.repeated) {
    return [];
  }
  if (field.default !== void 0) {
    return field.default;
  }
  switch (field.kind) {
    case "enum":
      return field.T.values[0].no;
    case "scalar":
      return scalarZeroValue(field.T, field.L);
    case "message":
      const T = field.T, value = new T();
      return T.fieldWrapper ? T.fieldWrapper.unwrapField(value) : value;
    case "map":
      throw "map fields are not allowed to be extensions";
  }
}
function filterUnknownFields(unknownFields, field) {
  if (!field.repeated && (field.kind == "enum" || field.kind == "scalar")) {
    for (let i = unknownFields.length - 1; i >= 0; --i) {
      if (unknownFields[i].no == field.no) {
        return [unknownFields[i]];
      }
    }
    return [];
  }
  return unknownFields.filter((uf) => uf.no === field.no);
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/proto-base64.js
var encTable = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/".split("");
var decTable = [];
for (let i = 0; i < encTable.length; i++)
  decTable[encTable[i].charCodeAt(0)] = i;
decTable["-".charCodeAt(0)] = encTable.indexOf("+");
decTable["_".charCodeAt(0)] = encTable.indexOf("/");
var protoBase64 = {
  /**
   * Decodes a base64 string to a byte array.
   *
   * - ignores white-space, including line breaks and tabs
   * - allows inner padding (can decode concatenated base64 strings)
   * - does not require padding
   * - understands base64url encoding:
   *   "-" instead of "+",
   *   "_" instead of "/",
   *   no padding
   */
  dec(base64Str) {
    let es = base64Str.length * 3 / 4;
    if (base64Str[base64Str.length - 2] == "=")
      es -= 2;
    else if (base64Str[base64Str.length - 1] == "=")
      es -= 1;
    let bytes = new Uint8Array(es), bytePos = 0, groupPos = 0, b, p = 0;
    for (let i = 0; i < base64Str.length; i++) {
      b = decTable[base64Str.charCodeAt(i)];
      if (b === void 0) {
        switch (base64Str[i]) {
          // @ts-ignore TS7029: Fallthrough case in switch
          case "=":
            groupPos = 0;
          // reset state when padding found
          // @ts-ignore TS7029: Fallthrough case in switch
          case "\n":
          case "\r":
          case "	":
          case " ":
            continue;
          // skip white-space, and padding
          default:
            throw Error("invalid base64 string.");
        }
      }
      switch (groupPos) {
        case 0:
          p = b;
          groupPos = 1;
          break;
        case 1:
          bytes[bytePos++] = p << 2 | (b & 48) >> 4;
          p = b;
          groupPos = 2;
          break;
        case 2:
          bytes[bytePos++] = (p & 15) << 4 | (b & 60) >> 2;
          p = b;
          groupPos = 3;
          break;
        case 3:
          bytes[bytePos++] = (p & 3) << 6 | b;
          groupPos = 0;
          break;
      }
    }
    if (groupPos == 1)
      throw Error("invalid base64 string.");
    return bytes.subarray(0, bytePos);
  },
  /**
   * Encode a byte array to a base64 string.
   */
  enc(bytes) {
    let base64 = "", groupPos = 0, b, p = 0;
    for (let i = 0; i < bytes.length; i++) {
      b = bytes[i];
      switch (groupPos) {
        case 0:
          base64 += encTable[b >> 2];
          p = (b & 3) << 4;
          groupPos = 1;
          break;
        case 1:
          base64 += encTable[p | b >> 4];
          p = (b & 15) << 2;
          groupPos = 2;
          break;
        case 2:
          base64 += encTable[p | b >> 6];
          base64 += encTable[b & 63];
          groupPos = 0;
          break;
      }
    }
    if (groupPos) {
      base64 += encTable[p];
      base64 += "=";
      if (groupPos == 1)
        base64 += "=";
    }
    return base64;
  }
};

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/extension-accessor.js
function getExtension(message, extension, options) {
  assertExtendee(extension, message);
  const opt = extension.runtime.bin.makeReadOptions(options);
  const ufs = filterUnknownFields(message.getType().runtime.bin.listUnknownFields(message), extension.field);
  const [container, get] = createExtensionContainer(extension);
  for (const uf of ufs) {
    extension.runtime.bin.readField(container, opt.readerFactory(uf.data), extension.field, uf.wireType, opt);
  }
  return get();
}
function setExtension(message, extension, value, options) {
  assertExtendee(extension, message);
  const readOpt = extension.runtime.bin.makeReadOptions(options);
  const writeOpt = extension.runtime.bin.makeWriteOptions(options);
  if (hasExtension(message, extension)) {
    const ufs = message.getType().runtime.bin.listUnknownFields(message).filter((uf) => uf.no != extension.field.no);
    message.getType().runtime.bin.discardUnknownFields(message);
    for (const uf of ufs) {
      message.getType().runtime.bin.onUnknownField(message, uf.no, uf.wireType, uf.data);
    }
  }
  const writer = writeOpt.writerFactory();
  let f = extension.field;
  if (!f.opt && !f.repeated && (f.kind == "enum" || f.kind == "scalar")) {
    f = Object.assign(Object.assign({}, extension.field), { opt: true });
  }
  extension.runtime.bin.writeField(f, value, writer, writeOpt);
  const reader = readOpt.readerFactory(writer.finish());
  while (reader.pos < reader.len) {
    const [no, wireType] = reader.tag();
    const data = reader.skip(wireType, no);
    message.getType().runtime.bin.onUnknownField(message, no, wireType, data);
  }
}
function hasExtension(message, extension) {
  const messageType = message.getType();
  return extension.extendee.typeName === messageType.typeName && !!messageType.runtime.bin.listUnknownFields(message).find((uf) => uf.no == extension.field.no);
}
function assertExtendee(extension, message) {
  assert(extension.extendee.typeName == message.getType().typeName, `extension ${extension.typeName} can only be applied to message ${extension.extendee.typeName}`);
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/reflect.js
function isFieldSet(field, target) {
  const localName = field.localName;
  if (field.repeated) {
    return target[localName].length > 0;
  }
  if (field.oneof) {
    return target[field.oneof.localName].case === localName;
  }
  switch (field.kind) {
    case "enum":
    case "scalar":
      if (field.opt || field.req) {
        return target[localName] !== void 0;
      }
      if (field.kind == "enum") {
        return target[localName] !== field.T.values[0].no;
      }
      return !isScalarZeroValue(field.T, target[localName]);
    case "message":
      return target[localName] !== void 0;
    case "map":
      return Object.keys(target[localName]).length > 0;
  }
}
function clearField(field, target) {
  const localName = field.localName;
  const implicitPresence = !field.opt && !field.req;
  if (field.repeated) {
    target[localName] = [];
  } else if (field.oneof) {
    target[field.oneof.localName] = { case: void 0 };
  } else {
    switch (field.kind) {
      case "map":
        target[localName] = {};
        break;
      case "enum":
        target[localName] = implicitPresence ? field.T.values[0].no : void 0;
        break;
      case "scalar":
        target[localName] = implicitPresence ? scalarZeroValue(field.T, field.L) : void 0;
        break;
      case "message":
        target[localName] = void 0;
        break;
    }
  }
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/is-message.js
function isMessage(arg, type) {
  if (arg === null || typeof arg != "object") {
    return false;
  }
  if (!Object.getOwnPropertyNames(Message.prototype).every((m) => m in arg && typeof arg[m] == "function")) {
    return false;
  }
  const actualType = arg.getType();
  if (actualType === null || typeof actualType != "function" || !("typeName" in actualType) || typeof actualType.typeName != "string") {
    return false;
  }
  return type === void 0 ? true : actualType.typeName == type.typeName;
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/field-wrapper.js
function wrapField(type, value) {
  if (isMessage(value) || !type.fieldWrapper) {
    return value;
  }
  return type.fieldWrapper.wrapField(value);
}
var wktWrapperToScalarType = {
  "google.protobuf.DoubleValue": ScalarType.DOUBLE,
  "google.protobuf.FloatValue": ScalarType.FLOAT,
  "google.protobuf.Int64Value": ScalarType.INT64,
  "google.protobuf.UInt64Value": ScalarType.UINT64,
  "google.protobuf.Int32Value": ScalarType.INT32,
  "google.protobuf.UInt32Value": ScalarType.UINT32,
  "google.protobuf.BoolValue": ScalarType.BOOL,
  "google.protobuf.StringValue": ScalarType.STRING,
  "google.protobuf.BytesValue": ScalarType.BYTES
};

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/json-format.js
var jsonReadDefaults = {
  ignoreUnknownFields: false
};
var jsonWriteDefaults = {
  emitDefaultValues: false,
  enumAsInteger: false,
  useProtoFieldName: false,
  prettySpaces: 0
};
function makeReadOptions(options) {
  return options ? Object.assign(Object.assign({}, jsonReadDefaults), options) : jsonReadDefaults;
}
function makeWriteOptions(options) {
  return options ? Object.assign(Object.assign({}, jsonWriteDefaults), options) : jsonWriteDefaults;
}
var tokenNull = Symbol();
var tokenIgnoredUnknownEnum = Symbol();
function makeJsonFormat() {
  return {
    makeReadOptions,
    makeWriteOptions,
    readMessage(type, json, options, message) {
      if (json == null || Array.isArray(json) || typeof json != "object") {
        throw new Error(`cannot decode message ${type.typeName} from JSON: ${debugJsonValue(json)}`);
      }
      message = message !== null && message !== void 0 ? message : new type();
      const oneofSeen = /* @__PURE__ */ new Map();
      const registry = options.typeRegistry;
      for (const [jsonKey, jsonValue] of Object.entries(json)) {
        const field = type.fields.findJsonName(jsonKey);
        if (field) {
          if (field.oneof) {
            if (jsonValue === null && field.kind == "scalar") {
              continue;
            }
            const seen = oneofSeen.get(field.oneof);
            if (seen !== void 0) {
              throw new Error(`cannot decode message ${type.typeName} from JSON: multiple keys for oneof "${field.oneof.name}" present: "${seen}", "${jsonKey}"`);
            }
            oneofSeen.set(field.oneof, jsonKey);
          }
          readField(message, jsonValue, field, options, type);
        } else {
          let found = false;
          if ((registry === null || registry === void 0 ? void 0 : registry.findExtension) && jsonKey.startsWith("[") && jsonKey.endsWith("]")) {
            const ext = registry.findExtension(jsonKey.substring(1, jsonKey.length - 1));
            if (ext && ext.extendee.typeName == type.typeName) {
              found = true;
              const [container, get] = createExtensionContainer(ext);
              readField(container, jsonValue, ext.field, options, ext);
              setExtension(message, ext, get(), options);
            }
          }
          if (!found && !options.ignoreUnknownFields) {
            throw new Error(`cannot decode message ${type.typeName} from JSON: key "${jsonKey}" is unknown`);
          }
        }
      }
      return message;
    },
    writeMessage(message, options) {
      const type = message.getType();
      const json = {};
      let field;
      try {
        for (field of type.fields.byNumber()) {
          if (!isFieldSet(field, message)) {
            if (field.req) {
              throw `required field not set`;
            }
            if (!options.emitDefaultValues) {
              continue;
            }
            if (!canEmitFieldDefaultValue(field)) {
              continue;
            }
          }
          const value = field.oneof ? message[field.oneof.localName].value : message[field.localName];
          const jsonValue = writeField(field, value, options);
          if (jsonValue !== void 0) {
            json[options.useProtoFieldName ? field.name : field.jsonName] = jsonValue;
          }
        }
        const registry = options.typeRegistry;
        if (registry === null || registry === void 0 ? void 0 : registry.findExtensionFor) {
          for (const uf of type.runtime.bin.listUnknownFields(message)) {
            const ext = registry.findExtensionFor(type.typeName, uf.no);
            if (ext && hasExtension(message, ext)) {
              const value = getExtension(message, ext, options);
              const jsonValue = writeField(ext.field, value, options);
              if (jsonValue !== void 0) {
                json[ext.field.jsonName] = jsonValue;
              }
            }
          }
        }
      } catch (e) {
        const m = field ? `cannot encode field ${type.typeName}.${field.name} to JSON` : `cannot encode message ${type.typeName} to JSON`;
        const r = e instanceof Error ? e.message : String(e);
        throw new Error(m + (r.length > 0 ? `: ${r}` : ""));
      }
      return json;
    },
    readScalar(type, json, longType) {
      return readScalar(type, json, longType !== null && longType !== void 0 ? longType : LongType.BIGINT, true);
    },
    writeScalar(type, value, emitDefaultValues) {
      if (value === void 0) {
        return void 0;
      }
      if (emitDefaultValues || isScalarZeroValue(type, value)) {
        return writeScalar(type, value);
      }
      return void 0;
    },
    debug: debugJsonValue
  };
}
function debugJsonValue(json) {
  if (json === null) {
    return "null";
  }
  switch (typeof json) {
    case "object":
      return Array.isArray(json) ? "array" : "object";
    case "string":
      return json.length > 100 ? "string" : `"${json.split('"').join('\\"')}"`;
    default:
      return String(json);
  }
}
function readField(target, jsonValue, field, options, parentType) {
  let localName = field.localName;
  if (field.repeated) {
    assert(field.kind != "map");
    if (jsonValue === null) {
      return;
    }
    if (!Array.isArray(jsonValue)) {
      throw new Error(`cannot decode field ${parentType.typeName}.${field.name} from JSON: ${debugJsonValue(jsonValue)}`);
    }
    const targetArray = target[localName];
    for (const jsonItem of jsonValue) {
      if (jsonItem === null) {
        throw new Error(`cannot decode field ${parentType.typeName}.${field.name} from JSON: ${debugJsonValue(jsonItem)}`);
      }
      switch (field.kind) {
        case "message":
          targetArray.push(field.T.fromJson(jsonItem, options));
          break;
        case "enum":
          const enumValue = readEnum(field.T, jsonItem, options.ignoreUnknownFields, true);
          if (enumValue !== tokenIgnoredUnknownEnum) {
            targetArray.push(enumValue);
          }
          break;
        case "scalar":
          try {
            targetArray.push(readScalar(field.T, jsonItem, field.L, true));
          } catch (e) {
            let m = `cannot decode field ${parentType.typeName}.${field.name} from JSON: ${debugJsonValue(jsonItem)}`;
            if (e instanceof Error && e.message.length > 0) {
              m += `: ${e.message}`;
            }
            throw new Error(m);
          }
          break;
      }
    }
  } else if (field.kind == "map") {
    if (jsonValue === null) {
      return;
    }
    if (typeof jsonValue != "object" || Array.isArray(jsonValue)) {
      throw new Error(`cannot decode field ${parentType.typeName}.${field.name} from JSON: ${debugJsonValue(jsonValue)}`);
    }
    const targetMap = target[localName];
    for (const [jsonMapKey, jsonMapValue] of Object.entries(jsonValue)) {
      if (jsonMapValue === null) {
        throw new Error(`cannot decode field ${parentType.typeName}.${field.name} from JSON: map value null`);
      }
      let key;
      try {
        key = readMapKey(field.K, jsonMapKey);
      } catch (e) {
        let m = `cannot decode map key for field ${parentType.typeName}.${field.name} from JSON: ${debugJsonValue(jsonValue)}`;
        if (e instanceof Error && e.message.length > 0) {
          m += `: ${e.message}`;
        }
        throw new Error(m);
      }
      switch (field.V.kind) {
        case "message":
          targetMap[key] = field.V.T.fromJson(jsonMapValue, options);
          break;
        case "enum":
          const enumValue = readEnum(field.V.T, jsonMapValue, options.ignoreUnknownFields, true);
          if (enumValue !== tokenIgnoredUnknownEnum) {
            targetMap[key] = enumValue;
          }
          break;
        case "scalar":
          try {
            targetMap[key] = readScalar(field.V.T, jsonMapValue, LongType.BIGINT, true);
          } catch (e) {
            let m = `cannot decode map value for field ${parentType.typeName}.${field.name} from JSON: ${debugJsonValue(jsonValue)}`;
            if (e instanceof Error && e.message.length > 0) {
              m += `: ${e.message}`;
            }
            throw new Error(m);
          }
          break;
      }
    }
  } else {
    if (field.oneof) {
      target = target[field.oneof.localName] = { case: localName };
      localName = "value";
    }
    switch (field.kind) {
      case "message":
        const messageType = field.T;
        if (jsonValue === null && messageType.typeName != "google.protobuf.Value") {
          return;
        }
        let currentValue = target[localName];
        if (isMessage(currentValue)) {
          currentValue.fromJson(jsonValue, options);
        } else {
          target[localName] = currentValue = messageType.fromJson(jsonValue, options);
          if (messageType.fieldWrapper && !field.oneof) {
            target[localName] = messageType.fieldWrapper.unwrapField(currentValue);
          }
        }
        break;
      case "enum":
        const enumValue = readEnum(field.T, jsonValue, options.ignoreUnknownFields, false);
        switch (enumValue) {
          case tokenNull:
            clearField(field, target);
            break;
          case tokenIgnoredUnknownEnum:
            break;
          default:
            target[localName] = enumValue;
            break;
        }
        break;
      case "scalar":
        try {
          const scalarValue = readScalar(field.T, jsonValue, field.L, false);
          switch (scalarValue) {
            case tokenNull:
              clearField(field, target);
              break;
            default:
              target[localName] = scalarValue;
              break;
          }
        } catch (e) {
          let m = `cannot decode field ${parentType.typeName}.${field.name} from JSON: ${debugJsonValue(jsonValue)}`;
          if (e instanceof Error && e.message.length > 0) {
            m += `: ${e.message}`;
          }
          throw new Error(m);
        }
        break;
    }
  }
}
function readMapKey(type, json) {
  if (type === ScalarType.BOOL) {
    switch (json) {
      case "true":
        json = true;
        break;
      case "false":
        json = false;
        break;
    }
  }
  return readScalar(type, json, LongType.BIGINT, true).toString();
}
function readScalar(type, json, longType, nullAsZeroValue) {
  if (json === null) {
    if (nullAsZeroValue) {
      return scalarZeroValue(type, longType);
    }
    return tokenNull;
  }
  switch (type) {
    // float, double: JSON value will be a number or one of the special string values "NaN", "Infinity", and "-Infinity".
    // Either numbers or strings are accepted. Exponent notation is also accepted.
    case ScalarType.DOUBLE:
    case ScalarType.FLOAT:
      if (json === "NaN")
        return Number.NaN;
      if (json === "Infinity")
        return Number.POSITIVE_INFINITY;
      if (json === "-Infinity")
        return Number.NEGATIVE_INFINITY;
      if (json === "") {
        break;
      }
      if (typeof json == "string" && json.trim().length !== json.length) {
        break;
      }
      if (typeof json != "string" && typeof json != "number") {
        break;
      }
      const float = Number(json);
      if (Number.isNaN(float)) {
        break;
      }
      if (!Number.isFinite(float)) {
        break;
      }
      if (type == ScalarType.FLOAT)
        assertFloat32(float);
      return float;
    // int32, fixed32, uint32: JSON value will be a decimal number. Either numbers or strings are accepted.
    case ScalarType.INT32:
    case ScalarType.FIXED32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32:
    case ScalarType.UINT32:
      let int32;
      if (typeof json == "number")
        int32 = json;
      else if (typeof json == "string" && json.length > 0) {
        if (json.trim().length === json.length)
          int32 = Number(json);
      }
      if (int32 === void 0)
        break;
      if (type == ScalarType.UINT32 || type == ScalarType.FIXED32)
        assertUInt32(int32);
      else
        assertInt32(int32);
      return int32;
    // int64, fixed64, uint64: JSON value will be a decimal string. Either numbers or strings are accepted.
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      if (typeof json != "number" && typeof json != "string")
        break;
      const long = protoInt64.parse(json);
      return longType ? long.toString() : long;
    case ScalarType.FIXED64:
    case ScalarType.UINT64:
      if (typeof json != "number" && typeof json != "string")
        break;
      const uLong = protoInt64.uParse(json);
      return longType ? uLong.toString() : uLong;
    // bool:
    case ScalarType.BOOL:
      if (typeof json !== "boolean")
        break;
      return json;
    // string:
    case ScalarType.STRING:
      if (typeof json !== "string") {
        break;
      }
      try {
        encodeURIComponent(json);
      } catch (e) {
        throw new Error("invalid UTF8");
      }
      return json;
    // bytes: JSON value will be the data encoded as a string using standard base64 encoding with paddings.
    // Either standard or URL-safe base64 encoding with/without paddings are accepted.
    case ScalarType.BYTES:
      if (json === "")
        return new Uint8Array(0);
      if (typeof json !== "string")
        break;
      return protoBase64.dec(json);
  }
  throw new Error();
}
function readEnum(type, json, ignoreUnknownFields, nullAsZeroValue) {
  if (json === null) {
    if (type.typeName == "google.protobuf.NullValue") {
      return 0;
    }
    return nullAsZeroValue ? type.values[0].no : tokenNull;
  }
  switch (typeof json) {
    case "number":
      if (Number.isInteger(json)) {
        return json;
      }
      break;
    case "string":
      const value = type.findName(json);
      if (value !== void 0) {
        return value.no;
      }
      if (ignoreUnknownFields) {
        return tokenIgnoredUnknownEnum;
      }
      break;
  }
  throw new Error(`cannot decode enum ${type.typeName} from JSON: ${debugJsonValue(json)}`);
}
function canEmitFieldDefaultValue(field) {
  if (field.repeated || field.kind == "map") {
    return true;
  }
  if (field.oneof) {
    return false;
  }
  if (field.kind == "message") {
    return false;
  }
  if (field.opt || field.req) {
    return false;
  }
  return true;
}
function writeField(field, value, options) {
  if (field.kind == "map") {
    assert(typeof value == "object" && value != null);
    const jsonObj = {};
    const entries = Object.entries(value);
    switch (field.V.kind) {
      case "scalar":
        for (const [entryKey, entryValue] of entries) {
          jsonObj[entryKey.toString()] = writeScalar(field.V.T, entryValue);
        }
        break;
      case "message":
        for (const [entryKey, entryValue] of entries) {
          jsonObj[entryKey.toString()] = entryValue.toJson(options);
        }
        break;
      case "enum":
        const enumType = field.V.T;
        for (const [entryKey, entryValue] of entries) {
          jsonObj[entryKey.toString()] = writeEnum(enumType, entryValue, options.enumAsInteger);
        }
        break;
    }
    return options.emitDefaultValues || entries.length > 0 ? jsonObj : void 0;
  }
  if (field.repeated) {
    assert(Array.isArray(value));
    const jsonArr = [];
    switch (field.kind) {
      case "scalar":
        for (let i = 0; i < value.length; i++) {
          jsonArr.push(writeScalar(field.T, value[i]));
        }
        break;
      case "enum":
        for (let i = 0; i < value.length; i++) {
          jsonArr.push(writeEnum(field.T, value[i], options.enumAsInteger));
        }
        break;
      case "message":
        for (let i = 0; i < value.length; i++) {
          jsonArr.push(value[i].toJson(options));
        }
        break;
    }
    return options.emitDefaultValues || jsonArr.length > 0 ? jsonArr : void 0;
  }
  switch (field.kind) {
    case "scalar":
      return writeScalar(field.T, value);
    case "enum":
      return writeEnum(field.T, value, options.enumAsInteger);
    case "message":
      return wrapField(field.T, value).toJson(options);
  }
}
function writeEnum(type, value, enumAsInteger) {
  var _a;
  assert(typeof value == "number");
  if (type.typeName == "google.protobuf.NullValue") {
    return null;
  }
  if (enumAsInteger) {
    return value;
  }
  const val = type.findNumber(value);
  return (_a = val === null || val === void 0 ? void 0 : val.name) !== null && _a !== void 0 ? _a : value;
}
function writeScalar(type, value) {
  switch (type) {
    // int32, fixed32, uint32: JSON value will be a decimal number. Either numbers or strings are accepted.
    case ScalarType.INT32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32:
    case ScalarType.FIXED32:
    case ScalarType.UINT32:
      assert(typeof value == "number");
      return value;
    // float, double: JSON value will be a number or one of the special string values "NaN", "Infinity", and "-Infinity".
    // Either numbers or strings are accepted. Exponent notation is also accepted.
    case ScalarType.FLOAT:
    // assertFloat32(value);
    case ScalarType.DOUBLE:
      assert(typeof value == "number");
      if (Number.isNaN(value))
        return "NaN";
      if (value === Number.POSITIVE_INFINITY)
        return "Infinity";
      if (value === Number.NEGATIVE_INFINITY)
        return "-Infinity";
      return value;
    // string:
    case ScalarType.STRING:
      assert(typeof value == "string");
      return value;
    // bool:
    case ScalarType.BOOL:
      assert(typeof value == "boolean");
      return value;
    // JSON value will be a decimal string. Either numbers or strings are accepted.
    case ScalarType.UINT64:
    case ScalarType.FIXED64:
    case ScalarType.INT64:
    case ScalarType.SFIXED64:
    case ScalarType.SINT64:
      assert(typeof value == "bigint" || typeof value == "string" || typeof value == "number");
      return value.toString();
    // bytes: JSON value will be the data encoded as a string using standard base64 encoding with paddings.
    // Either standard or URL-safe base64 encoding with/without paddings are accepted.
    case ScalarType.BYTES:
      assert(value instanceof Uint8Array);
      return protoBase64.enc(value);
  }
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/binary-format.js
var unknownFieldsSymbol = Symbol("@bufbuild/protobuf/unknown-fields");
var readDefaults = {
  readUnknownFields: true,
  readerFactory: (bytes) => new BinaryReader(bytes)
};
var writeDefaults = {
  writeUnknownFields: true,
  writerFactory: () => new BinaryWriter()
};
function makeReadOptions2(options) {
  return options ? Object.assign(Object.assign({}, readDefaults), options) : readDefaults;
}
function makeWriteOptions2(options) {
  return options ? Object.assign(Object.assign({}, writeDefaults), options) : writeDefaults;
}
function makeBinaryFormat() {
  return {
    makeReadOptions: makeReadOptions2,
    makeWriteOptions: makeWriteOptions2,
    listUnknownFields(message) {
      var _a;
      return (_a = message[unknownFieldsSymbol]) !== null && _a !== void 0 ? _a : [];
    },
    discardUnknownFields(message) {
      delete message[unknownFieldsSymbol];
    },
    writeUnknownFields(message, writer) {
      const m = message;
      const c = m[unknownFieldsSymbol];
      if (c) {
        for (const f of c) {
          writer.tag(f.no, f.wireType).raw(f.data);
        }
      }
    },
    onUnknownField(message, no, wireType, data) {
      const m = message;
      if (!Array.isArray(m[unknownFieldsSymbol])) {
        m[unknownFieldsSymbol] = [];
      }
      m[unknownFieldsSymbol].push({ no, wireType, data });
    },
    readMessage(message, reader, lengthOrEndTagFieldNo, options, delimitedMessageEncoding) {
      const type = message.getType();
      const end = delimitedMessageEncoding ? reader.len : reader.pos + lengthOrEndTagFieldNo;
      let fieldNo, wireType;
      while (reader.pos < end) {
        [fieldNo, wireType] = reader.tag();
        if (delimitedMessageEncoding === true && wireType == WireType.EndGroup) {
          break;
        }
        const field = type.fields.find(fieldNo);
        if (!field) {
          const data = reader.skip(wireType, fieldNo);
          if (options.readUnknownFields) {
            this.onUnknownField(message, fieldNo, wireType, data);
          }
          continue;
        }
        readField2(message, reader, field, wireType, options);
      }
      if (delimitedMessageEncoding && // eslint-disable-line @typescript-eslint/strict-boolean-expressions
      (wireType != WireType.EndGroup || fieldNo !== lengthOrEndTagFieldNo)) {
        throw new Error(`invalid end group tag`);
      }
    },
    readField: readField2,
    writeMessage(message, writer, options) {
      const type = message.getType();
      for (const field of type.fields.byNumber()) {
        if (!isFieldSet(field, message)) {
          if (field.req) {
            throw new Error(`cannot encode field ${type.typeName}.${field.name} to binary: required field not set`);
          }
          continue;
        }
        const value = field.oneof ? message[field.oneof.localName].value : message[field.localName];
        writeField2(field, value, writer, options);
      }
      if (options.writeUnknownFields) {
        this.writeUnknownFields(message, writer);
      }
      return writer;
    },
    writeField(field, value, writer, options) {
      if (value === void 0) {
        return void 0;
      }
      writeField2(field, value, writer, options);
    }
  };
}
function readField2(target, reader, field, wireType, options) {
  let { repeated, localName } = field;
  if (field.oneof) {
    target = target[field.oneof.localName];
    if (target.case != localName) {
      delete target.value;
    }
    target.case = localName;
    localName = "value";
  }
  switch (field.kind) {
    case "scalar":
    case "enum":
      const scalarType = field.kind == "enum" ? ScalarType.INT32 : field.T;
      let read = readScalar2;
      if (field.kind == "scalar" && field.L > 0) {
        read = readScalarLTString;
      }
      if (repeated) {
        let arr = target[localName];
        const isPacked = wireType == WireType.LengthDelimited && scalarType != ScalarType.STRING && scalarType != ScalarType.BYTES;
        if (isPacked) {
          let e = reader.uint32() + reader.pos;
          while (reader.pos < e) {
            arr.push(read(reader, scalarType));
          }
        } else {
          arr.push(read(reader, scalarType));
        }
      } else {
        target[localName] = read(reader, scalarType);
      }
      break;
    case "message":
      const messageType = field.T;
      if (repeated) {
        target[localName].push(readMessageField(reader, new messageType(), options, field));
      } else {
        if (isMessage(target[localName])) {
          readMessageField(reader, target[localName], options, field);
        } else {
          target[localName] = readMessageField(reader, new messageType(), options, field);
          if (messageType.fieldWrapper && !field.oneof && !field.repeated) {
            target[localName] = messageType.fieldWrapper.unwrapField(target[localName]);
          }
        }
      }
      break;
    case "map":
      let [mapKey, mapVal] = readMapEntry(field, reader, options);
      target[localName][mapKey] = mapVal;
      break;
  }
}
function readMessageField(reader, message, options, field) {
  const format = message.getType().runtime.bin;
  const delimited = field === null || field === void 0 ? void 0 : field.delimited;
  format.readMessage(
    message,
    reader,
    delimited ? field.no : reader.uint32(),
    // eslint-disable-line @typescript-eslint/strict-boolean-expressions
    options,
    delimited
  );
  return message;
}
function readMapEntry(field, reader, options) {
  const length = reader.uint32(), end = reader.pos + length;
  let key, val;
  while (reader.pos < end) {
    const [fieldNo] = reader.tag();
    switch (fieldNo) {
      case 1:
        key = readScalar2(reader, field.K);
        break;
      case 2:
        switch (field.V.kind) {
          case "scalar":
            val = readScalar2(reader, field.V.T);
            break;
          case "enum":
            val = reader.int32();
            break;
          case "message":
            val = readMessageField(reader, new field.V.T(), options, void 0);
            break;
        }
        break;
    }
  }
  if (key === void 0) {
    key = scalarZeroValue(field.K, LongType.BIGINT);
  }
  if (typeof key != "string" && typeof key != "number") {
    key = key.toString();
  }
  if (val === void 0) {
    switch (field.V.kind) {
      case "scalar":
        val = scalarZeroValue(field.V.T, LongType.BIGINT);
        break;
      case "enum":
        val = field.V.T.values[0].no;
        break;
      case "message":
        val = new field.V.T();
        break;
    }
  }
  return [key, val];
}
function readScalarLTString(reader, type) {
  const v = readScalar2(reader, type);
  return typeof v == "bigint" ? v.toString() : v;
}
function readScalar2(reader, type) {
  switch (type) {
    case ScalarType.STRING:
      return reader.string();
    case ScalarType.BOOL:
      return reader.bool();
    case ScalarType.DOUBLE:
      return reader.double();
    case ScalarType.FLOAT:
      return reader.float();
    case ScalarType.INT32:
      return reader.int32();
    case ScalarType.INT64:
      return reader.int64();
    case ScalarType.UINT64:
      return reader.uint64();
    case ScalarType.FIXED64:
      return reader.fixed64();
    case ScalarType.BYTES:
      return reader.bytes();
    case ScalarType.FIXED32:
      return reader.fixed32();
    case ScalarType.SFIXED32:
      return reader.sfixed32();
    case ScalarType.SFIXED64:
      return reader.sfixed64();
    case ScalarType.SINT64:
      return reader.sint64();
    case ScalarType.UINT32:
      return reader.uint32();
    case ScalarType.SINT32:
      return reader.sint32();
  }
}
function writeField2(field, value, writer, options) {
  assert(value !== void 0);
  const repeated = field.repeated;
  switch (field.kind) {
    case "scalar":
    case "enum":
      let scalarType = field.kind == "enum" ? ScalarType.INT32 : field.T;
      if (repeated) {
        assert(Array.isArray(value));
        if (field.packed) {
          writePacked(writer, scalarType, field.no, value);
        } else {
          for (const item of value) {
            writeScalar2(writer, scalarType, field.no, item);
          }
        }
      } else {
        writeScalar2(writer, scalarType, field.no, value);
      }
      break;
    case "message":
      if (repeated) {
        assert(Array.isArray(value));
        for (const item of value) {
          writeMessageField(writer, options, field, item);
        }
      } else {
        writeMessageField(writer, options, field, value);
      }
      break;
    case "map":
      assert(typeof value == "object" && value != null);
      for (const [key, val] of Object.entries(value)) {
        writeMapEntry(writer, options, field, key, val);
      }
      break;
  }
}
function writeMapEntry(writer, options, field, key, value) {
  writer.tag(field.no, WireType.LengthDelimited);
  writer.fork();
  let keyValue = key;
  switch (field.K) {
    case ScalarType.INT32:
    case ScalarType.FIXED32:
    case ScalarType.UINT32:
    case ScalarType.SFIXED32:
    case ScalarType.SINT32:
      keyValue = Number.parseInt(key);
      break;
    case ScalarType.BOOL:
      assert(key == "true" || key == "false");
      keyValue = key == "true";
      break;
  }
  writeScalar2(writer, field.K, 1, keyValue);
  switch (field.V.kind) {
    case "scalar":
      writeScalar2(writer, field.V.T, 2, value);
      break;
    case "enum":
      writeScalar2(writer, ScalarType.INT32, 2, value);
      break;
    case "message":
      assert(value !== void 0);
      writer.tag(2, WireType.LengthDelimited).bytes(value.toBinary(options));
      break;
  }
  writer.join();
}
function writeMessageField(writer, options, field, value) {
  const message = wrapField(field.T, value);
  if (field.delimited)
    writer.tag(field.no, WireType.StartGroup).raw(message.toBinary(options)).tag(field.no, WireType.EndGroup);
  else
    writer.tag(field.no, WireType.LengthDelimited).bytes(message.toBinary(options));
}
function writeScalar2(writer, type, fieldNo, value) {
  assert(value !== void 0);
  let [wireType, method] = scalarTypeInfo(type);
  writer.tag(fieldNo, wireType)[method](value);
}
function writePacked(writer, type, fieldNo, value) {
  if (!value.length) {
    return;
  }
  writer.tag(fieldNo, WireType.LengthDelimited).fork();
  let [, method] = scalarTypeInfo(type);
  for (let i = 0; i < value.length; i++) {
    writer[method](value[i]);
  }
  writer.join();
}
function scalarTypeInfo(type) {
  let wireType = WireType.Varint;
  switch (type) {
    case ScalarType.BYTES:
    case ScalarType.STRING:
      wireType = WireType.LengthDelimited;
      break;
    case ScalarType.DOUBLE:
    case ScalarType.FIXED64:
    case ScalarType.SFIXED64:
      wireType = WireType.Bit64;
      break;
    case ScalarType.FIXED32:
    case ScalarType.SFIXED32:
    case ScalarType.FLOAT:
      wireType = WireType.Bit32;
      break;
  }
  const method = ScalarType[type].toLowerCase();
  return [wireType, method];
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/util-common.js
function makeUtilCommon() {
  return {
    setEnumType,
    initPartial(source, target) {
      if (source === void 0) {
        return;
      }
      const type = target.getType();
      for (const member of type.fields.byMember()) {
        const localName = member.localName, t = target, s = source;
        if (s[localName] == null) {
          continue;
        }
        switch (member.kind) {
          case "oneof":
            const sk = s[localName].case;
            if (sk === void 0) {
              continue;
            }
            const sourceField = member.findField(sk);
            let val = s[localName].value;
            if (sourceField && sourceField.kind == "message" && !isMessage(val, sourceField.T)) {
              val = new sourceField.T(val);
            } else if (sourceField && sourceField.kind === "scalar" && sourceField.T === ScalarType.BYTES) {
              val = toU8Arr(val);
            }
            t[localName] = { case: sk, value: val };
            break;
          case "scalar":
          case "enum":
            let copy = s[localName];
            if (member.T === ScalarType.BYTES) {
              copy = member.repeated ? copy.map(toU8Arr) : toU8Arr(copy);
            }
            t[localName] = copy;
            break;
          case "map":
            switch (member.V.kind) {
              case "scalar":
              case "enum":
                if (member.V.T === ScalarType.BYTES) {
                  for (const [k, v] of Object.entries(s[localName])) {
                    t[localName][k] = toU8Arr(v);
                  }
                } else {
                  Object.assign(t[localName], s[localName]);
                }
                break;
              case "message":
                const messageType = member.V.T;
                for (const k of Object.keys(s[localName])) {
                  let val2 = s[localName][k];
                  if (!messageType.fieldWrapper) {
                    val2 = new messageType(val2);
                  }
                  t[localName][k] = val2;
                }
                break;
            }
            break;
          case "message":
            const mt = member.T;
            if (member.repeated) {
              t[localName] = s[localName].map((val2) => isMessage(val2, mt) ? val2 : new mt(val2));
            } else {
              const val2 = s[localName];
              if (mt.fieldWrapper) {
                if (
                  // We can't use BytesValue.typeName as that will create a circular import
                  mt.typeName === "google.protobuf.BytesValue"
                ) {
                  t[localName] = toU8Arr(val2);
                } else {
                  t[localName] = val2;
                }
              } else {
                t[localName] = isMessage(val2, mt) ? val2 : new mt(val2);
              }
            }
            break;
        }
      }
    },
    // TODO use isFieldSet() here to support future field presence
    equals(type, a, b) {
      if (a === b) {
        return true;
      }
      if (!a || !b) {
        return false;
      }
      return type.fields.byMember().every((m) => {
        const va = a[m.localName];
        const vb = b[m.localName];
        if (m.repeated) {
          if (va.length !== vb.length) {
            return false;
          }
          switch (m.kind) {
            case "message":
              return va.every((a2, i) => m.T.equals(a2, vb[i]));
            case "scalar":
              return va.every((a2, i) => scalarEquals(m.T, a2, vb[i]));
            case "enum":
              return va.every((a2, i) => scalarEquals(ScalarType.INT32, a2, vb[i]));
          }
          throw new Error(`repeated cannot contain ${m.kind}`);
        }
        switch (m.kind) {
          case "message":
            return m.T.equals(va, vb);
          case "enum":
            return scalarEquals(ScalarType.INT32, va, vb);
          case "scalar":
            return scalarEquals(m.T, va, vb);
          case "oneof":
            if (va.case !== vb.case) {
              return false;
            }
            const s = m.findField(va.case);
            if (s === void 0) {
              return true;
            }
            switch (s.kind) {
              case "message":
                return s.T.equals(va.value, vb.value);
              case "enum":
                return scalarEquals(ScalarType.INT32, va.value, vb.value);
              case "scalar":
                return scalarEquals(s.T, va.value, vb.value);
            }
            throw new Error(`oneof cannot contain ${s.kind}`);
          case "map":
            const keys = Object.keys(va).concat(Object.keys(vb));
            switch (m.V.kind) {
              case "message":
                const messageType = m.V.T;
                return keys.every((k) => messageType.equals(va[k], vb[k]));
              case "enum":
                return keys.every((k) => scalarEquals(ScalarType.INT32, va[k], vb[k]));
              case "scalar":
                const scalarType = m.V.T;
                return keys.every((k) => scalarEquals(scalarType, va[k], vb[k]));
            }
            break;
        }
      });
    },
    // TODO use isFieldSet() here to support future field presence
    clone(message) {
      const type = message.getType(), target = new type(), any = target;
      for (const member of type.fields.byMember()) {
        const source = message[member.localName];
        let copy;
        if (member.repeated) {
          copy = source.map(cloneSingularField);
        } else if (member.kind == "map") {
          copy = any[member.localName];
          for (const [key, v] of Object.entries(source)) {
            copy[key] = cloneSingularField(v);
          }
        } else if (member.kind == "oneof") {
          const f = member.findField(source.case);
          copy = f ? { case: source.case, value: cloneSingularField(source.value) } : { case: void 0 };
        } else {
          copy = cloneSingularField(source);
        }
        any[member.localName] = copy;
      }
      for (const uf of type.runtime.bin.listUnknownFields(message)) {
        type.runtime.bin.onUnknownField(any, uf.no, uf.wireType, uf.data);
      }
      return target;
    }
  };
}
function cloneSingularField(value) {
  if (value === void 0) {
    return value;
  }
  if (isMessage(value)) {
    return value.clone();
  }
  if (value instanceof Uint8Array) {
    const c = new Uint8Array(value.byteLength);
    c.set(value);
    return c;
  }
  return value;
}
function toU8Arr(input) {
  return input instanceof Uint8Array ? input : new Uint8Array(input);
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/proto-runtime.js
function makeProtoRuntime(syntax, newFieldList, initFields) {
  return {
    syntax,
    json: makeJsonFormat(),
    bin: makeBinaryFormat(),
    util: Object.assign(Object.assign({}, makeUtilCommon()), {
      newFieldList,
      initFields
    }),
    makeMessageType(typeName, fields, opt) {
      return makeMessageType(this, typeName, fields, opt);
    },
    makeEnum,
    makeEnumType,
    getEnumType,
    makeExtension(typeName, extendee, field) {
      return makeExtension(this, typeName, extendee, field);
    }
  };
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/field-list.js
var InternalFieldList = class {
  constructor(fields, normalizer) {
    this._fields = fields;
    this._normalizer = normalizer;
  }
  findJsonName(jsonName) {
    if (!this.jsonNames) {
      const t = {};
      for (const f of this.list()) {
        t[f.jsonName] = t[f.name] = f;
      }
      this.jsonNames = t;
    }
    return this.jsonNames[jsonName];
  }
  find(fieldNo) {
    if (!this.numbers) {
      const t = {};
      for (const f of this.list()) {
        t[f.no] = f;
      }
      this.numbers = t;
    }
    return this.numbers[fieldNo];
  }
  list() {
    if (!this.all) {
      this.all = this._normalizer(this._fields);
    }
    return this.all;
  }
  byNumber() {
    if (!this.numbersAsc) {
      this.numbersAsc = this.list().concat().sort((a, b) => a.no - b.no);
    }
    return this.numbersAsc;
  }
  byMember() {
    if (!this.members) {
      this.members = [];
      const a = this.members;
      let o;
      for (const f of this.list()) {
        if (f.oneof) {
          if (f.oneof !== o) {
            o = f.oneof;
            a.push(o);
          }
        } else {
          a.push(f);
        }
      }
    }
    return this.members;
  }
};

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/names.js
function localFieldName(protoName, inOneof) {
  const name = protoCamelCase(protoName);
  if (inOneof) {
    return name;
  }
  return safeObjectProperty(safeMessageProperty(name));
}
function localOneofName(protoName) {
  return localFieldName(protoName, false);
}
var fieldJsonName = protoCamelCase;
function protoCamelCase(snakeCase) {
  let capNext = false;
  const b = [];
  for (let i = 0; i < snakeCase.length; i++) {
    let c = snakeCase.charAt(i);
    switch (c) {
      case "_":
        capNext = true;
        break;
      case "0":
      case "1":
      case "2":
      case "3":
      case "4":
      case "5":
      case "6":
      case "7":
      case "8":
      case "9":
        b.push(c);
        capNext = false;
        break;
      default:
        if (capNext) {
          capNext = false;
          c = c.toUpperCase();
        }
        b.push(c);
        break;
    }
  }
  return b.join("");
}
var reservedObjectProperties = /* @__PURE__ */ new Set([
  // names reserved by JavaScript
  "constructor",
  "toString",
  "toJSON",
  "valueOf"
]);
var reservedMessageProperties = /* @__PURE__ */ new Set([
  // names reserved by the runtime
  "getType",
  "clone",
  "equals",
  "fromBinary",
  "fromJson",
  "fromJsonString",
  "toBinary",
  "toJson",
  "toJsonString",
  // names reserved by the runtime for the future
  "toObject"
]);
var fallback = (name) => `${name}$`;
var safeMessageProperty = (name) => {
  if (reservedMessageProperties.has(name)) {
    return fallback(name);
  }
  return name;
};
var safeObjectProperty = (name) => {
  if (reservedObjectProperties.has(name)) {
    return fallback(name);
  }
  return name;
};

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/field.js
var InternalOneofInfo = class {
  constructor(name) {
    this.kind = "oneof";
    this.repeated = false;
    this.packed = false;
    this.opt = false;
    this.req = false;
    this.default = void 0;
    this.fields = [];
    this.name = name;
    this.localName = localOneofName(name);
  }
  addField(field) {
    assert(field.oneof === this, `field ${field.name} not one of ${this.name}`);
    this.fields.push(field);
  }
  findField(localName) {
    if (!this._lookup) {
      this._lookup = /* @__PURE__ */ Object.create(null);
      for (let i = 0; i < this.fields.length; i++) {
        this._lookup[this.fields[i].localName] = this.fields[i];
      }
    }
    return this._lookup[localName];
  }
};

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/private/field-normalize.js
function normalizeFieldInfos(fieldInfos, packedByDefault) {
  var _a, _b, _c, _d, _e, _f;
  const r = [];
  let o;
  for (const field of typeof fieldInfos == "function" ? fieldInfos() : fieldInfos) {
    const f = field;
    f.localName = localFieldName(field.name, field.oneof !== void 0);
    f.jsonName = (_a = field.jsonName) !== null && _a !== void 0 ? _a : fieldJsonName(field.name);
    f.repeated = (_b = field.repeated) !== null && _b !== void 0 ? _b : false;
    if (field.kind == "scalar") {
      f.L = (_c = field.L) !== null && _c !== void 0 ? _c : LongType.BIGINT;
    }
    f.delimited = (_d = field.delimited) !== null && _d !== void 0 ? _d : false;
    f.req = (_e = field.req) !== null && _e !== void 0 ? _e : false;
    f.opt = (_f = field.opt) !== null && _f !== void 0 ? _f : false;
    if (field.packed === void 0) {
      if (packedByDefault) {
        f.packed = field.kind == "enum" || field.kind == "scalar" && field.T != ScalarType.BYTES && field.T != ScalarType.STRING;
      } else {
        f.packed = false;
      }
    }
    if (field.oneof !== void 0) {
      const ooname = typeof field.oneof == "string" ? field.oneof : field.oneof.name;
      if (!o || o.name != ooname) {
        o = new InternalOneofInfo(ooname);
      }
      f.oneof = o;
      o.addField(f);
    }
    r.push(f);
  }
  return r;
}

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/proto3.js
var proto3 = makeProtoRuntime(
  "proto3",
  (fields) => {
    return new InternalFieldList(fields, (source) => normalizeFieldInfos(source, true));
  },
  // TODO merge with proto2 and initExtensionField, also see initPartial, equals, clone
  (target) => {
    for (const member of target.getType().fields.byMember()) {
      if (member.opt) {
        continue;
      }
      const name = member.localName, t = target;
      if (member.repeated) {
        t[name] = [];
        continue;
      }
      switch (member.kind) {
        case "oneof":
          t[name] = { case: void 0 };
          break;
        case "enum":
          t[name] = 0;
          break;
        case "map":
          t[name] = {};
          break;
        case "scalar":
          t[name] = scalarZeroValue(member.T, member.L);
          break;
        case "message":
          break;
      }
    }
  }
);

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/service-type.js
var MethodKind;
(function(MethodKind2) {
  MethodKind2[MethodKind2["Unary"] = 0] = "Unary";
  MethodKind2[MethodKind2["ServerStreaming"] = 1] = "ServerStreaming";
  MethodKind2[MethodKind2["ClientStreaming"] = 2] = "ClientStreaming";
  MethodKind2[MethodKind2["BiDiStreaming"] = 3] = "BiDiStreaming";
})(MethodKind || (MethodKind = {}));
var MethodIdempotency;
(function(MethodIdempotency2) {
  MethodIdempotency2[MethodIdempotency2["NoSideEffects"] = 1] = "NoSideEffects";
  MethodIdempotency2[MethodIdempotency2["Idempotent"] = 2] = "Idempotent";
})(MethodIdempotency || (MethodIdempotency = {}));

// node_modules/.pnpm/@bufbuild+protobuf@1.10.0/node_modules/@bufbuild/protobuf/dist/esm/google/protobuf/empty_pb.js
var Empty = class _Empty extends Message {
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static fromBinary(bytes, options) {
    return new _Empty().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _Empty().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _Empty().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_Empty, a, b);
  }
};
Empty.runtime = proto3;
Empty.typeName = "google.protobuf.Empty";
Empty.fields = proto3.util.newFieldList(() => []);

// src/gen/clientserver_pb.ts
var Card = class _Card extends Message {
  /**
   * @generated from field: string card_id = 1;
   */
  cardId = "";
  /**
   * @generated from field: string name = 2;
   */
  name = "";
  /**
   * @generated from field: string code = 3;
   */
  code = "";
  /**
   * @generated from field: string description = 4;
   */
  description = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.Card";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "card_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "name",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 3,
      name: "code",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 4,
      name: "description",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _Card().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _Card().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _Card().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_Card, a, b);
  }
};
var RetrieveCardsDataResponse = class _RetrieveCardsDataResponse extends Message {
  /**
   * @generated from field: repeated com.sweetloveinyourheart.kittens.clients.Card cards = 1;
   */
  cards = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.RetrieveCardsDataResponse";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "cards", kind: "message", T: Card, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _RetrieveCardsDataResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _RetrieveCardsDataResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _RetrieveCardsDataResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_RetrieveCardsDataResponse, a, b);
  }
};
var User = class _User extends Message {
  /**
   * @generated from field: string user_id = 1;
   */
  userId = "";
  /**
   * @generated from field: string username = 2;
   */
  username = "";
  /**
   * @generated from field: string full_name = 3;
   */
  fullName = "";
  /**
   * @generated from field: int32 status = 4;
   */
  status = 0;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.User";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "user_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "username",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 3,
      name: "full_name",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 4,
      name: "status",
      kind: "scalar",
      T: 5
      /* ScalarType.INT32 */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _User().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _User().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _User().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_User, a, b);
  }
};
var CreateNewGuestUserRequest = class _CreateNewGuestUserRequest extends Message {
  /**
   * Required: Username of the guest user
   *
   * @generated from field: string username = 1;
   */
  username = "";
  /**
   * Required: Full name of the guest user
   *
   * @generated from field: string full_name = 2;
   */
  fullName = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "username",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "full_name",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _CreateNewGuestUserRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _CreateNewGuestUserRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _CreateNewGuestUserRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_CreateNewGuestUserRequest, a, b);
  }
};
var CreateNewGuestUserResponse = class _CreateNewGuestUserResponse extends Message {
  /**
   * The user basic info
   *
   * @generated from field: com.sweetloveinyourheart.kittens.clients.User user = 1;
   */
  user;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.CreateNewGuestUserResponse";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "user", kind: "message", T: User }
  ]);
  static fromBinary(bytes, options) {
    return new _CreateNewGuestUserResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _CreateNewGuestUserResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _CreateNewGuestUserResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_CreateNewGuestUserResponse, a, b);
  }
};
var GuestLoginRequest = class _GuestLoginRequest extends Message {
  /**
   * Required: UUID of the guest user
   *
   * @generated from field: string user_id = 1;
   */
  userId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GuestLoginRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "user_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _GuestLoginRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GuestLoginRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GuestLoginRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GuestLoginRequest, a, b);
  }
};
var GuestLoginResponse = class _GuestLoginResponse extends Message {
  /**
   * The user basic info
   *
   * @generated from field: com.sweetloveinyourheart.kittens.clients.User user = 1;
   */
  user;
  /**
   * The session token for this user.
   *
   * @generated from field: string token = 2;
   */
  token = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GuestLoginResponse";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "user", kind: "message", T: User },
    {
      no: 2,
      name: "token",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _GuestLoginResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GuestLoginResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GuestLoginResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GuestLoginResponse, a, b);
  }
};
var UserProfileResponse = class _UserProfileResponse extends Message {
  /**
   * @generated from field: com.sweetloveinyourheart.kittens.clients.User user = 1;
   */
  user;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.UserProfileResponse";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "user", kind: "message", T: User }
  ]);
  static fromBinary(bytes, options) {
    return new _UserProfileResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _UserProfileResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _UserProfileResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_UserProfileResponse, a, b);
  }
};
var PlayersProfileRequest = class _PlayersProfileRequest extends Message {
  /**
   * Required: UUID of the guest user
   *
   * @generated from field: repeated string user_ids = 1;
   */
  userIds = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.PlayersProfileRequest";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "user_ids", kind: "scalar", T: 9, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _PlayersProfileRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _PlayersProfileRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _PlayersProfileRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_PlayersProfileRequest, a, b);
  }
};
var PlayersProfileResponse = class _PlayersProfileResponse extends Message {
  /**
   * @generated from field: repeated com.sweetloveinyourheart.kittens.clients.User users = 1;
   */
  users = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.PlayersProfileResponse";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "users", kind: "message", T: User, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _PlayersProfileResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _PlayersProfileResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _PlayersProfileResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_PlayersProfileResponse, a, b);
  }
};
var Lobby = class _Lobby extends Message {
  /**
   * @generated from field: string lobby_id = 1;
   */
  lobbyId = "";
  /**
   * @generated from field: string lobby_code = 2;
   */
  lobbyCode = "";
  /**
   * @generated from field: string lobby_name = 3;
   */
  lobbyName = "";
  /**
   * @generated from field: string host_user_id = 4;
   */
  hostUserId = "";
  /**
   * @generated from field: repeated string participants = 5;
   */
  participants = [];
  /**
   * @generated from field: optional string match_id = 6;
   */
  matchId;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.Lobby";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "lobby_code",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 3,
      name: "lobby_name",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 4,
      name: "host_user_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    { no: 5, name: "participants", kind: "scalar", T: 9, repeated: true },
    { no: 6, name: "match_id", kind: "scalar", T: 9, opt: true }
  ]);
  static fromBinary(bytes, options) {
    return new _Lobby().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _Lobby().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _Lobby().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_Lobby, a, b);
  }
};
var CreateLobbyRequest = class _CreateLobbyRequest extends Message {
  /**
   * @generated from field: string lobby_name = 1;
   */
  lobbyName = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.CreateLobbyRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_name",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _CreateLobbyRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _CreateLobbyRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _CreateLobbyRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_CreateLobbyRequest, a, b);
  }
};
var CreateLobbyResponse = class _CreateLobbyResponse extends Message {
  /**
   * @generated from field: string lobby_id = 1;
   */
  lobbyId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.CreateLobbyResponse";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _CreateLobbyResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _CreateLobbyResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _CreateLobbyResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_CreateLobbyResponse, a, b);
  }
};
var GetLobbyRequest = class _GetLobbyRequest extends Message {
  /**
   * @generated from field: string lobby_id = 1;
   */
  lobbyId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GetLobbyRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _GetLobbyRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GetLobbyRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GetLobbyRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GetLobbyRequest, a, b);
  }
};
var GetLobbyReply = class _GetLobbyReply extends Message {
  /**
   * @generated from field: com.sweetloveinyourheart.kittens.clients.Lobby lobby = 1;
   */
  lobby;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GetLobbyReply";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "lobby", kind: "message", T: Lobby }
  ]);
  static fromBinary(bytes, options) {
    return new _GetLobbyReply().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GetLobbyReply().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GetLobbyReply().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GetLobbyReply, a, b);
  }
};
var JoinLobbyRequest = class _JoinLobbyRequest extends Message {
  /**
   * @generated from field: string lobby_code = 1;
   */
  lobbyCode = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.JoinLobbyRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_code",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _JoinLobbyRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _JoinLobbyRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _JoinLobbyRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_JoinLobbyRequest, a, b);
  }
};
var JoinLobbyResponse = class _JoinLobbyResponse extends Message {
  /**
   * @generated from field: string lobby_id = 1;
   */
  lobbyId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.JoinLobbyResponse";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _JoinLobbyResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _JoinLobbyResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _JoinLobbyResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_JoinLobbyResponse, a, b);
  }
};
var LeaveLobbyRequest = class _LeaveLobbyRequest extends Message {
  /**
   * @generated from field: string lobby_id = 1;
   */
  lobbyId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.LeaveLobbyRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _LeaveLobbyRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _LeaveLobbyRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _LeaveLobbyRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_LeaveLobbyRequest, a, b);
  }
};
var LeaveLobbyResponse = class _LeaveLobbyResponse extends Message {
  /**
   * @generated from field: string lobby_id = 1;
   */
  lobbyId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.LeaveLobbyResponse";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _LeaveLobbyResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _LeaveLobbyResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _LeaveLobbyResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_LeaveLobbyResponse, a, b);
  }
};
var StartMatchRequest = class _StartMatchRequest extends Message {
  /**
   * @generated from field: string lobby_id = 1;
   */
  lobbyId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.StartMatchRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "lobby_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _StartMatchRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _StartMatchRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _StartMatchRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_StartMatchRequest, a, b);
  }
};
var Game = class _Game extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  /**
   * @generated from field: com.sweetloveinyourheart.kittens.clients.Game.Phase game_phase = 2;
   */
  gamePhase = 0 /* INITIALIZING */;
  /**
   * @generated from field: string player_turn = 3;
   */
  playerTurn = "";
  /**
   * @generated from field: repeated com.sweetloveinyourheart.kittens.clients.Game.Player players = 4;
   */
  players = [];
  /**
   * @generated from field: map<string, com.sweetloveinyourheart.kittens.clients.Game.PlayerHand> player_hands = 5;
   */
  playerHands = {};
  /**
   * @generated from field: com.sweetloveinyourheart.kittens.clients.Game.Desk desk = 6;
   */
  desk;
  /**
   * @generated from field: string executing_action = 7;
   */
  executingAction = "";
  /**
   * @generated from field: string affected_player = 8;
   */
  affectedPlayer = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.Game";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    { no: 2, name: "game_phase", kind: "enum", T: proto3.getEnumType(Game_Phase) },
    {
      no: 3,
      name: "player_turn",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    { no: 4, name: "players", kind: "message", T: Game_Player, repeated: true },
    { no: 5, name: "player_hands", kind: "map", K: 9, V: { kind: "message", T: Game_PlayerHand } },
    { no: 6, name: "desk", kind: "message", T: Game_Desk },
    {
      no: 7,
      name: "executing_action",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 8,
      name: "affected_player",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _Game().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _Game().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _Game().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_Game, a, b);
  }
};
var Game_Phase = /* @__PURE__ */ ((Game_Phase2) => {
  Game_Phase2[Game_Phase2["INITIALIZING"] = 0] = "INITIALIZING";
  Game_Phase2[Game_Phase2["TURN_START"] = 1] = "TURN_START";
  Game_Phase2[Game_Phase2["ACTION_PHASE"] = 2] = "ACTION_PHASE";
  Game_Phase2[Game_Phase2["CARD_DRAWING"] = 3] = "CARD_DRAWING";
  Game_Phase2[Game_Phase2["TURN_END"] = 4] = "TURN_END";
  Game_Phase2[Game_Phase2["GAME_OVER"] = 5] = "GAME_OVER";
  return Game_Phase2;
})(Game_Phase || {});
proto3.util.setEnumType(Game_Phase, "com.sweetloveinyourheart.kittens.clients.Game.Phase", [
  { no: 0, name: "INITIALIZING" },
  { no: 1, name: "TURN_START" },
  { no: 2, name: "ACTION_PHASE" },
  { no: 3, name: "CARD_DRAWING" },
  { no: 4, name: "TURN_END" },
  { no: 5, name: "GAME_OVER" }
]);
var Game_Player = class _Game_Player extends Message {
  /**
   * @generated from field: string player_id = 1;
   */
  playerId = "";
  /**
   * @generated from field: bool active = 2;
   */
  active = false;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.Game.Player";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "player_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "active",
      kind: "scalar",
      T: 8
      /* ScalarType.BOOL */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _Game_Player().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _Game_Player().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _Game_Player().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_Game_Player, a, b);
  }
};
var Game_PlayerHand = class _Game_PlayerHand extends Message {
  /**
   * @generated from field: int32 remaining_cards = 1;
   */
  remainingCards = 0;
  /**
   * @generated from field: repeated string hands = 2;
   */
  hands = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.Game.PlayerHand";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "remaining_cards",
      kind: "scalar",
      T: 5
      /* ScalarType.INT32 */
    },
    { no: 2, name: "hands", kind: "scalar", T: 9, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _Game_PlayerHand().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _Game_PlayerHand().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _Game_PlayerHand().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_Game_PlayerHand, a, b);
  }
};
var Game_Desk = class _Game_Desk extends Message {
  /**
   * @generated from field: string desk_id = 1;
   */
  deskId = "";
  /**
   * @generated from field: int32 remaining_cards = 2;
   */
  remainingCards = 0;
  /**
   * @generated from field: repeated string discard_pile = 3;
   */
  discardPile = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.Game.Desk";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "desk_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "remaining_cards",
      kind: "scalar",
      T: 5
      /* ScalarType.INT32 */
    },
    { no: 3, name: "discard_pile", kind: "scalar", T: 9, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _Game_Desk().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _Game_Desk().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _Game_Desk().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_Game_Desk, a, b);
  }
};
var GameMetaData = class _GameMetaData extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  /**
   * @generated from field: repeated string players = 2;
   */
  players = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GameMetaData";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    { no: 2, name: "players", kind: "scalar", T: 9, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _GameMetaData().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GameMetaData().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GameMetaData().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GameMetaData, a, b);
  }
};
var StreamGameRequest = class _StreamGameRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.StreamGameRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _StreamGameRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _StreamGameRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _StreamGameRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_StreamGameRequest, a, b);
  }
};
var StreamGameReply = class _StreamGameReply extends Message {
  /**
   * @generated from field: com.sweetloveinyourheart.kittens.clients.Game game_state = 1;
   */
  gameState;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.StreamGameReply";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "game_state", kind: "message", T: Game }
  ]);
  static fromBinary(bytes, options) {
    return new _StreamGameReply().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _StreamGameReply().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _StreamGameReply().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_StreamGameReply, a, b);
  }
};
var GetGameMetaDataRequest = class _GetGameMetaDataRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GetGameMetaDataRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _GetGameMetaDataRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GetGameMetaDataRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GetGameMetaDataRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GetGameMetaDataRequest, a, b);
  }
};
var GetGameMetaDataResponse = class _GetGameMetaDataResponse extends Message {
  /**
   * @generated from field: com.sweetloveinyourheart.kittens.clients.GameMetaData meta = 1;
   */
  meta;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GetGameMetaDataResponse";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "meta", kind: "message", T: GameMetaData }
  ]);
  static fromBinary(bytes, options) {
    return new _GetGameMetaDataResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GetGameMetaDataResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GetGameMetaDataResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GetGameMetaDataResponse, a, b);
  }
};
var PlayCardsRequest = class _PlayCardsRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  /**
   * @generated from field: repeated string card_ids = 2;
   */
  cardIds = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.PlayCardsRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    { no: 2, name: "card_ids", kind: "scalar", T: 9, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _PlayCardsRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _PlayCardsRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _PlayCardsRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_PlayCardsRequest, a, b);
  }
};
var PeekCardsRequest = class _PeekCardsRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  /**
   * @generated from field: string desk_id = 2;
   */
  deskId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.PeekCardsRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "desk_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _PeekCardsRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _PeekCardsRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _PeekCardsRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_PeekCardsRequest, a, b);
  }
};
var PeekCardsResponse = class _PeekCardsResponse extends Message {
  /**
   * @generated from field: repeated string card_ids = 1;
   */
  cardIds = [];
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.PeekCardsResponse";
  static fields = proto3.util.newFieldList(() => [
    { no: 1, name: "card_ids", kind: "scalar", T: 9, repeated: true }
  ]);
  static fromBinary(bytes, options) {
    return new _PeekCardsResponse().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _PeekCardsResponse().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _PeekCardsResponse().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_PeekCardsResponse, a, b);
  }
};
var SelectAffectedPlayerRequest = class _SelectAffectedPlayerRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  /**
   * @generated from field: string player_id = 2;
   */
  playerId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.SelectAffectedPlayerRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "player_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _SelectAffectedPlayerRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _SelectAffectedPlayerRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _SelectAffectedPlayerRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_SelectAffectedPlayerRequest, a, b);
  }
};
var StealCardRequest = class _StealCardRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  /**
   * @generated from field: optional string card_id = 2;
   */
  cardId;
  /**
   * @generated from field: optional int32 card_index = 3;
   */
  cardIndex;
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.StealCardRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    { no: 2, name: "card_id", kind: "scalar", T: 9, opt: true },
    { no: 3, name: "card_index", kind: "scalar", T: 5, opt: true }
  ]);
  static fromBinary(bytes, options) {
    return new _StealCardRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _StealCardRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _StealCardRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_StealCardRequest, a, b);
  }
};
var GiveCardRequest = class _GiveCardRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  /**
   * @generated from field: string card_id = 2;
   */
  cardId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.GiveCardRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    },
    {
      no: 2,
      name: "card_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _GiveCardRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _GiveCardRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _GiveCardRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_GiveCardRequest, a, b);
  }
};
var DrawCardsRequest = class _DrawCardsRequest extends Message {
  /**
   * @generated from field: string game_id = 1;
   */
  gameId = "";
  constructor(data) {
    super();
    proto3.util.initPartial(data, this);
  }
  static runtime = proto3;
  static typeName = "com.sweetloveinyourheart.kittens.clients.DrawCardsRequest";
  static fields = proto3.util.newFieldList(() => [
    {
      no: 1,
      name: "game_id",
      kind: "scalar",
      T: 9
      /* ScalarType.STRING */
    }
  ]);
  static fromBinary(bytes, options) {
    return new _DrawCardsRequest().fromBinary(bytes, options);
  }
  static fromJson(jsonValue, options) {
    return new _DrawCardsRequest().fromJson(jsonValue, options);
  }
  static fromJsonString(jsonString, options) {
    return new _DrawCardsRequest().fromJsonString(jsonString, options);
  }
  static equals(a, b) {
    return proto3.util.equals(_DrawCardsRequest, a, b);
  }
};

// src/gen/clientserver_connect.ts
var ClientServer = {
  typeName: "com.sweetloveinyourheart.kittens.clients.ClientServer",
  methods: {
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.RetrieveCardsData
     */
    retrieveCardsData: {
      name: "RetrieveCardsData",
      I: Empty,
      O: RetrieveCardsDataResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.CreateNewGuestUser
     */
    createNewGuestUser: {
      name: "CreateNewGuestUser",
      I: CreateNewGuestUserRequest,
      O: CreateNewGuestUserResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GuestLogin
     */
    guestLogin: {
      name: "GuestLogin",
      I: GuestLoginRequest,
      O: GuestLoginResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetUserProfile
     */
    getUserProfile: {
      name: "GetUserProfile",
      I: Empty,
      O: UserProfileResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetPlayersProfile
     */
    getPlayersProfile: {
      name: "GetPlayersProfile",
      I: PlayersProfileRequest,
      O: PlayersProfileResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.CreateLobby
     */
    createLobby: {
      name: "CreateLobby",
      I: CreateLobbyRequest,
      O: CreateLobbyResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetLobby
     */
    getLobby: {
      name: "GetLobby",
      I: GetLobbyRequest,
      O: GetLobbyReply,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StreamLobby
     */
    streamLobby: {
      name: "StreamLobby",
      I: GetLobbyRequest,
      O: GetLobbyReply,
      kind: MethodKind.ServerStreaming
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.JoinLobby
     */
    joinLobby: {
      name: "JoinLobby",
      I: JoinLobbyRequest,
      O: JoinLobbyResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.LeaveLobby
     */
    leaveLobby: {
      name: "LeaveLobby",
      I: LeaveLobbyRequest,
      O: LeaveLobbyResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StartMatch
     */
    startMatch: {
      name: "StartMatch",
      I: StartMatchRequest,
      O: Empty,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GetGameMetaData
     */
    getGameMetaData: {
      name: "GetGameMetaData",
      I: GetGameMetaDataRequest,
      O: GetGameMetaDataResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StreamGame
     */
    streamGame: {
      name: "StreamGame",
      I: StreamGameRequest,
      O: StreamGameReply,
      kind: MethodKind.ServerStreaming
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.PlayCards
     */
    playCards: {
      name: "PlayCards",
      I: PlayCardsRequest,
      O: Empty,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.PeekCards
     */
    peekCards: {
      name: "PeekCards",
      I: PeekCardsRequest,
      O: PeekCardsResponse,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.DrawCards
     */
    drawCards: {
      name: "DrawCards",
      I: DrawCardsRequest,
      O: Empty,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.SelectAffectedPlayer
     */
    selectAffectedPlayer: {
      name: "SelectAffectedPlayer",
      I: SelectAffectedPlayerRequest,
      O: Empty,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.StealCard
     */
    stealCard: {
      name: "StealCard",
      I: StealCardRequest,
      O: Empty,
      kind: MethodKind.Unary
    },
    /**
     * @generated from rpc com.sweetloveinyourheart.kittens.clients.ClientServer.GiveCard
     */
    giveCard: {
      name: "GiveCard",
      I: GiveCardRequest,
      O: Empty,
      kind: MethodKind.Unary
    }
  }
};
export {
  Card,
  ClientServer,
  CreateLobbyRequest,
  CreateLobbyResponse,
  CreateNewGuestUserRequest,
  CreateNewGuestUserResponse,
  DrawCardsRequest,
  Game,
  GameMetaData,
  Game_Desk,
  Game_Phase,
  Game_Player,
  Game_PlayerHand,
  GetGameMetaDataRequest,
  GetGameMetaDataResponse,
  GetLobbyReply,
  GetLobbyRequest,
  GiveCardRequest,
  GuestLoginRequest,
  GuestLoginResponse,
  JoinLobbyRequest,
  JoinLobbyResponse,
  LeaveLobbyRequest,
  LeaveLobbyResponse,
  Lobby,
  PeekCardsRequest,
  PeekCardsResponse,
  PlayCardsRequest,
  PlayersProfileRequest,
  PlayersProfileResponse,
  RetrieveCardsDataResponse,
  SelectAffectedPlayerRequest,
  StartMatchRequest,
  StealCardRequest,
  StreamGameReply,
  StreamGameRequest,
  User,
  UserProfileResponse
};
