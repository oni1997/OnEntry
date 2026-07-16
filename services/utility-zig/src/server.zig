const std = @import("std");
const net = std.net;

pub const Server = struct {
    allocator: std.mem.Allocator,
    address: net.Address,
    listener: net.Server,
    routes: std.StringHashMap(*const fn (*Context) anyerror!void),

    pub const Context = struct {
        request: Request,
        writer: std.io.AnyWriter,
        allocator: std.mem.Allocator,

        pub fn json(self: *Context, value: anytype) !void {
            const json_str = try std.json.stringifyAlloc(self.allocator, value, .{});
            defer self.allocator.free(json_str);
            _ = try self.writer.print("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: {d}\r\n\r\n{s}", .{ json_str.len, json_str });
        }
    };

    pub const Request = struct {
        method: []const u8,
        path: []const u8,
        body: []const u8,
    };

    pub fn init(allocator: std.mem.Allocator, port: u16) !Server {
        const address = try net.Address.initIp4(.{ 0, 0, 0, 0 }, port);
        var server = Server{
            .allocator = allocator,
            .address = address,
            .listener = net.Server.init(.{ .reuse_address = true }),
            .routes = std.StringHashMap(*const fn (*Context) anyerror!void).init(allocator),
        };
        try server.listener.listen(address);
        return server;
    }

    pub fn deinit(self: *Server) void {
        self.listener.deinit();
        self.routes.deinit();
    }

    pub fn get(self: *Server, path: []const u8, handler: *const fn (*Context) anyerror!void) !void {
        const owned = try self.allocator.dupe(u8, path);
        try self.routes.put(owned, handler);
    }

    pub fn post(self: *Server, path: []const u8, handler: *const fn (*Context) anyerror!void) !void {
        const owned = try self.allocator.dupe(u8, path);
        try self.routes.put(owned, handler);
    }

    pub fn run(self: *Server) !void {
        std.log.info("Utility service listening on {any}", .{self.address});
        while (true) {
            const conn = self.listener.accept() catch |err| {
                std.log.warn("Accept error: {any}", .{err});
                continue;
            };
            _ = try std.Thread.spawn(.{}, handleConnection, .{self, conn});
        }
    }

    fn handleConnection(self: *Server, conn: net.Server.Connection) void {
        defer conn.stream.close() catch {};
        const reader = conn.stream.reader();
        const writer = conn.stream.writer();

        var buf: [4096]u8 = undefined;
        const read = reader.read(&buf) catch 0;
        if (read == 0) return;

        const request_str = buf[0..read];
        var lines = std.mem.splitSequence(u8, request_str, "\r\n");
        const request_line = lines.next() orelse return;
        var parts = std.mem.splitSequence(u8, request_line, " ");
        const method = parts.next() orelse "";
        const path_str = parts.next() orelse "";
        const route_path = std.mem.trimLeft(u8, path_str, "/");

        std.log.info("{s} {s}", .{ method, route_path });

        if (self.routes.get(route_path)) |handler| {
            const ctx = Context{
                .request = Request{
                    .method = method,
                    .path = route_path,
                    .body = request_str,
                },
                .writer = writer.any(),
                .allocator = self.allocator,
            };
            handler(&ctx) catch |err| {
                std.log.err("Handler error: {any}", .{err});
                _ = writer.print("HTTP/1.1 500 Internal Server Error\r\n\r\n") catch {};
            };
        } else {
            _ = writer.print("HTTP/1.1 404 Not Found\r\n\r\n") catch {};
        }
    }
};
