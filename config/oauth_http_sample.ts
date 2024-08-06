
// server.ts
import { serve } from "bun";

// 型定義
interface RequestHandler {
  fetch(req: Request): Promise<Response>;
}

// サーバーの作成
const server: RequestHandler = {
  async fetch(req: Request): Promise<Response> {
    const url = new URL(req.url);

    // OAuth2リダイレクトを処理
    if (url.pathname === "/auth/twitch/callback") {
      const hash = url.hash;
      return new Response(`
        <!DOCTYPE html>
        <html>
          <head>
            <meta charset="UTF-8" />
          </head>
          <body>
            <script>
              if (location.hash) {
                alert(\`トークン情報は${location.hash}だよ！\`);
              }
            </script>
          </body>
        </html>
      `, { headers: { "Content-Type": "text/html" } });
    }

    // メインHTMLページを提供
    return new Response(`
      <!DOCTYPE html>
      <html>
        <head>
          <meta charset="UTF-8" />
        </head>
        <body>
          <a href="https://id.twitch.tv/oauth2/authorize?client_id=&redirect_uri=http://localhost:3000&response_type=token&scope=bits:read+channel:bot+channel:read:hype_train+channel:read:polls+channel:read:predictions+channel:read:redemptions+channel:manage:redemptions+channel:read:subscriptions+moderation:read+moderator:manage:announcements+moderator:manage:automod+moderator:read:chat_messages+moderator:manage:chat_messages+moderator:read:chat_settings+moderator:read:chatters+moderator:read:followers+moderator:read:shoutouts+user:bot+user:edit+user:edit:broadcast+user:read:chat+user:manage:chat_color+user:read:emotes+user:read:follows+user:read:moderated_channels+user:write:chat" >
            Twitch認証
          </a>
        </body>
      </html>
    `, { headers: { "Content-Type": "text/html" } });
  }
};

// サーバーの起動
serve({
  port: 3000,
  fetch: server.fetch,
});

console.log("Server is running on http://localhost:3000");
