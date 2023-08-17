type WsMessagePayload = {
  Type: string;
  Reference: string;
  Channel: string;
  Payload: string;
};

export type CloseCallback = (event: CloseEvent) => void;
export type OpenCallback = (event: Event) => void;
export type MessageCallback = (
  playload: WsMessagePayload,
  event: MessageEvent<WsMessagePayload>
) => void;

type CanalMessageCallbacks = Record<string, MessageCallback[] | undefined>;

export class Ws {
  private conn: WebSocket;
  private onCloseCallbacks: CloseCallback[] = [];
  private onOpenCallbacks: OpenCallback[] = [];
  private onAnyMessageCallbacks: MessageCallback[] = [];
  private onCanalMessageCallbacks: CanalMessageCallbacks = {};

  constructor(url: URL, headers: Record<string, string> = {}) {
    if (typeof window === "undefined") {
      throw Error("Ws can only be use in browser");
    }

    if (!("WebSocket" in window)) {
      throw Error("The browser does not support WebSockets.");
    }

    const params = new URLSearchParams(headers).toString();
    const queries = params.length > 0 ? `?${params}` : "";
    const host = url.origin.replace("http", "ws") + "/ws" + queries;

    // Create WebSocket Connection
    this.conn = new WebSocket(host);

    this.conn.onopen = (event) => {
      this.onOpenCallbacks.forEach((cb) => cb(event));
    };

    this.conn.onclose = (event) => {
      this.onCloseCallbacks.forEach((cb) => cb(event));
    };

    // On Message
    this.conn.onmessage = (event: MessageEvent<WsMessagePayload>) => {
      const channel = event.data?.Channel;
      const reference = event.data?.Reference;
      // Any message listeners
      this.onAnyMessageCallbacks.forEach((cb) => cb(event.data, event));

      [reference, channel].forEach((canal) => {
        if (canal) {
          this.onCanalMessageCallbacks[canal]?.forEach((cb) =>
            cb(event.data, event)
          );
        }
      });
    };
  }

  private _onCanalMessage(canal: string, cb: MessageCallback) {
    let callbacks = this.onCanalMessageCallbacks[canal];
    if (!callbacks) {
      callbacks = [cb];
      this.onCanalMessageCallbacks[canal] = callbacks;
    } else {
      callbacks.push(cb);
    }

    return () => {
      this.onCanalMessageCallbacks[canal] = callbacks?.filter(
        (func) => func !== cb
      );
    };
  }

  public onChannelMessage(channel: string, cb: MessageCallback) {
    return this._onCanalMessage(channel, cb);
  }

  public onReferenceMessage(reference: string, cb: MessageCallback) {
    return this._onCanalMessage(reference, cb);
  }

  public onAnyMessage(cb: MessageCallback) {
    this.onAnyMessageCallbacks.push(cb);
    return () => {
      this.onAnyMessageCallbacks = this.onAnyMessageCallbacks.filter(
        (func) => func !== cb
      );
    };
  }

  public onOpen(cb: OpenCallback) {
    this.onOpenCallbacks.push(cb);
    return () => {
      this.onOpenCallbacks = this.onOpenCallbacks.filter((func) => func !== cb);
    };
  }

  public onClose(cb: CloseCallback) {
    this.onCloseCallbacks.push(cb);
    return () => {
      this.onCloseCallbacks = this.onCloseCallbacks.filter(
        (func) => func !== cb
      );
    };
  }

  public close() {
    this.conn.close();
  }
}
