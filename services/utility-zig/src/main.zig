const std = @import("std");
const api = @import("api.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    std.log.info("OnEntry Utility Service starting on port 8084", .{});

    try api.startServer(allocator, 8084);
}
