const API_BASE = process.env.API_BASE_URL || 'http://api-server:8080';

type User = { id: number; name: string };

export default async function Home() {
  const res = await fetch(`${API_BASE}/users`, { cache: 'no-store' });
  const users = await res.json();

  return (
    <main>
      <h1>Users</h1>
      <ul>
        {users.map((user: User) => (
          <li key={user.id}>{user.name}</li>
        ))}
      </ul>
    </main>
  );
}
