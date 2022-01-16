import List "mo:base/List";
import Iter "mo:base/Iter";
import Principal "mo:base/Principal";
import Time "mo:base/Time";

actor {
    // public type Message = Text;
    public type Message = {
        text: Text;
        time: Time.Time;
    };

    public type Microblog = actor {
        follow: shared(Principal) -> async (); // 添加关注对象
        follows: shared query () -> async [Principal]; // 返回关注列表
        post: shared (Text) -> async (); // 发布新消息
        allPosts: shared query () -> async [Message]; // 返回所有发布的消息
        allTimeline: shared () -> async [Message]; // 返回所有关注对象发布的消息
        posts: shared query (since: Time.Time) -> async [Message]; // 返回指定时间之后发布的消息
        timeline: shared (since: Time.Time) -> async [Message]; // 返回所有关注对象指定时间之后发布的消息
    };

    // stable 修饰变量，使得 canister 在升级之后，变量内容不丢失
    stable var followed : List.List<Principal> = List.nil();

    public shared func follow(id: Principal) : async() {
        followed := List.push(id, followed);
    };

    public shared query func follows() : async [Principal] {
        List.toArray(followed);
    };

    stable var messages : List.List<Message> = List.nil();

    public shared (msg) func post(text: Text) : async () {
        //assert(Principal.toText(msg.caller) == "ljxyp-mxm6h-xj3cu-vglnb-mbjqm-vybgi-djdt2-lqgqj-tttmw-dm7mx-tqe");
        var message = {
            text = text;
            time = Time.now();
        };
        messages := List.push(message, messages);
    };

    public shared query func allPosts() : async [Message] {
        List.toArray(messages);
    };
    
    public shared query func posts(since: Time.Time) : async [Message] {
        var resultMsg : List.List<Message> = List.nil();

        for (message in Iter.fromList(messages)) {
            if (message.time <= since) {
                resultMsg := List.push(message, resultMsg);
            };
        };

        List.toArray(resultMsg);
    };

    public shared func allTimeline() : async [Message] {
        var all : List.List<Message> = List.nil();

        for (id in Iter.fromList(followed)) {
            let canister : Microblog = actor(Principal.toText(id));
            let msgs = await canister.allPosts();
            for (msg in Iter.fromArray(msgs)) {
                all := List.push(msg, all);
            };
        };

        List.toArray(all);
    };

    public shared func timeline(since: Time.Time) : async [Message] {
        var all : List.List<Message> = List.nil();

        for (id in Iter.fromList(followed)) {
            let canister : Microblog = actor(Principal.toText(id));
            let msgs = await canister.posts(since);
            for (msg in Iter.fromArray(msgs)) {
                all := List.push(msg, all);
            };
        };

        List.toArray(all);
    };
};
