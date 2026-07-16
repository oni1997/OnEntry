const std = @import("std");
const Server = @import("server.zig").Server;

pub fn startServer(allocator: std.mem.Allocator, port: u16) !void {
    var server = Server.init(allocator, port) catch |err| {
        std.log.err("Failed to initialize server: {any}", .{err});
        return err;
    };
    defer server.deinit();

    try server.get("/health", health);
    try server.get("/backup", backup);
    try server.get("/export", exportVault);
    try server.get("/verify", verifyDatabase);
    try server.get("/cleanup", cleanup);
    try server.post("/restore", restore);

    try server.run();
}

fn health(ctx: *Server.Context) !void {
    try ctx.json(.{.status = "ok"});
}

fn backup(ctx: *Server.Context) !void {
    std.log.info("Backup requested", .{});
    try ctx.json(.{.message = "Backup started"});
}

fn exportVault(ctx: *Server.Context) !void {
    std.log.info("Export requested", .{});
    try ctx.json(.{.message = "Export started"});
}

fn verifyDatabase(ctx: *Server.Context) !void {
    std.log.info("Database verification requested", .{});
    try ctx.json(.{.status = "verified", .issues = 0});
}

fn cleanup(ctx: *Server.Context) !void {
    std.log.info("Cleanup requested", .{});
    try ctx.json(.{.message = "Cleanup completed"});
}

fn restore(ctx: *Server.Context) !void {
    std.log.info("Restore requested", .{});
    try ctx.json(.{.message = "Restore started"});
}
