import List "mo:base/List";
import Iter "mo:base/Iter";
import Principal "mo:base/Principal";
import Time "mo:base/Time";

actor {
    // public type Message = Text;
    public type Message = {
        author: Text;
        text: Text;
        time: Time.Time;
    };

    public type Microblog = actor {
        follow: shared(Principal) -> async (); // 添加关注对象
        follows: shared query () -> async [Principal]; // 返回关注列表
        post: shared (Text) -> async (); // 发布新消息
        //posts: shared query () -> async [Message]; // 返回所有发布的消息
        //timeline: shared () -> async [Message]; // 返回所有关注对象发布的消息
        posts: shared query (since: Time.Time) -> async [Message]; // 返回指定时间之后发布的消息
        timeline: shared (since: Time.Time) -> async [Message]; // 返回所有关注对象指定时间之后发布的消息
        get_name: shared query () -> async ?Text;
    };

    // stable 修饰变量，使得 canister 在升级之后，变量内容不丢失
    stable var author: Text = "";
    
    public shared func set_name(name: Text) {
        author := name;
    };

    public shared func get_name() : async ?Text {
        return ?author;
    };
    
    stable var followed : List.List<Principal> = List.nil();

    public shared func follow(id: Principal) : async() {
        followed := List.push(id, followed);
    };

    public shared query func follows() : async [Principal] {
        List.toArray(followed);
    };

    stable var messages : List.List<Message> = List.nil();

    public shared (msg) func post(opt: Text, text: Text) : async () {
        //assert(Principal.toText(msg.caller) == "ljxyp-mxm6h-xj3cu-vglnb-mbjqm-vybgi-djdt2-lqgqj-tttmw-dm7mx-tqe");
        assert(opt == author);
        var message = {
            author = opt;
            text = text;
            time = Time.now();
        };
        messages := List.push(message, messages);
    };

    /*public shared query func posts() : async [Message] {
        List.toArray(messages);
    };*/
    
    public shared query func posts(since: Time.Time) : async [Message] {
        var resultMsg : List.List<Message> = List.nil();

        for (message in Iter.fromList(messages)) {
            if (message.time <= since) {
                resultMsg := List.push(message, resultMsg);
            };
        };

        List.toArray(resultMsg);
    };

    /*public shared func timeline() : async [Message] {
        var all : List.List<Message> = List.nil();

        for (id in Iter.fromList(followed)) {
            let canister : Microblog = actor(Principal.toText(id));
            let msgs = await canister.posts();
            for (msg in Iter.fromArray(msgs)) {
                all := List.push(msg, all);
            };
        };

        List.toArray(all);
    };*/

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

    // 获取关注者列表
    public shared func allFollowings() : async [Text] {
        var all : List.List<Text> = List.nil();

        for (id in Iter.fromList(followed)) {
            let canister : Microblog = actor(Principal.toText(id));
            let optionalName = await canister.get_name();
            switch (optionalName) {
                case (null) {  };
                case (?name) { all := List.push(name, all); };
            }
        };

        List.toArray(all);
    };

    // 获取关注者的消息列表
    public shared func followingMessages(canisterID: Text) : async [Message] {
        var all : List.List<Message> = List.nil();

        let canister : Microblog = actor(canisterID);
        let msgs = await canister.posts(Time.now());
        for (msg in Iter.fromArray(msgs)) {
            all := List.push(msg, all);
        };

        List.toArray(all);
    };
};

