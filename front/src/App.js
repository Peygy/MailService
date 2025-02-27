import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate, Navigate } from "react-router-dom";
import axios from "axios";

const API_URL = "http://localhost:8081/api/v1";

function App() {
  const [auth, setAuth] = useState(() => {
    const storedAuth = localStorage.getItem('auth');
    return storedAuth ? storedAuth : null;  // Храним только строку, а не объект
  });

  const login = (username, password) => {
    const encoded = btoa(`${username}:${password}`);  // Закодированная строка
    setAuth(encoded);
    localStorage.setItem('auth', encoded);  // Сохраняем только строку
  };

  const authHeaders = auth ? { Authorization: `Basic ${auth}` } : {};

  return (
    <Router>
      <nav>
        <Link to="/login">Login</Link>
        <Link to="/register">Register</Link>
        <Link to="/inbox">Inbox</Link>
        <Link to="/sent">Sent</Link>
        <Link to="/send">SendMail</Link>
        <Link to="/admin">Admin</Link>
      </nav>
      <Routes>
        <Route path="/login" element={<Login onLogin={login} />} />
        <Route path="/register" element={<Register />} />
        <Route path="/inbox" element={auth ? <Inbox authHeaders={authHeaders} /> : <Navigate to="/login" />} />
        <Route path="/sent" element={auth ? <Sent authHeaders={authHeaders} /> : <Navigate to="/login" />} />
        <Route path="/send" element={auth ? <SendMail authHeaders={authHeaders} /> : <Navigate to="/login" />} />
        <Route path="/admin" element={auth ? <Admin authHeaders={authHeaders} /> : <Navigate to="/login" />} />
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </Router>
  );
}

function Login({ onLogin }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      await axios.post(`${API_URL}/login`, { email, password });
      onLogin(email, password);  // Сохраняем Basic Auth
      navigate("/inbox");
    } catch (err) {
      alert("Login failed");
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button onClick={handleLogin}>Login</button>
    </div>
  );
}

function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleRegister = async () => {
    try {
      await axios.post(`${API_URL}/register`, { email, password });
      alert("Registration successful");
    } catch (err) {
      alert("Registration failed");
    }
  };

  return (
    <div>
      <h2>Register</h2>
      <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button onClick={handleRegister}>Register</button>
    </div>
  );
}

function Inbox({ authHeaders }) {
  const [mails, setMails] = useState([]);

  useEffect(() => {
    axios.get(`${API_URL}/mail/inbox`, { headers: authHeaders })
      .then((res) => setMails(res.data.mails))
      .catch(() => alert("Failed to load inbox"));
  }, [authHeaders]);

  return (
    <div>
      <h2>Inbox</h2>
      {mails.length === 0 ? (
        <p>No emails in inbox</p>
      ) : (
        <table border="1" cellPadding="5">
          <thead>
            <tr>
              <th>Sender</th>
              <th>Subject</th>
              <th>Body</th>
              <th>Received At</th>
            </tr>
          </thead>
          <tbody>
            {mails.map((mail) => (
              <tr key={mail.ID}>
                <td>{mail.Sender}</td>
                <td>{mail.Subject}</td>
                <td>{mail.Body}</td>
                <td>{new Date(mail.CreatedAt).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

function Sent({ authHeaders }) {
  const [mails, setMails] = useState([]);

  useEffect(() => {
    axios.get(`${API_URL}/mail/sent`, { headers: authHeaders })
      .then((res) => setMails(res.data.mails))
      .catch(() => alert("Failed to load sent mails"));
  }, [authHeaders]);

  return (
    <div>
      <h2>Sent Mails</h2>
      {mails.length === 0 ? (
        <p>No sent emails</p>
      ) : (
        <table border="1" cellPadding="5">
          <thead>
            <tr>
              <th>Receiver</th>
              <th>Subject</th>
              <th>Body</th>
              <th>Sent At</th>
            </tr>
          </thead>
          <tbody>
            {mails.map((mail) => (
              <tr key={mail.ID}>
                <td>{mail.Receiver}</td>
                <td>{mail.Subject}</td>
                <td>{mail.Body}</td>
                <td>{new Date(mail.CreatedAt).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}


function SendMail({ authHeaders }) {
  const [receiver, setReceiver] = useState("");
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");

  const handleSend = async () => {
    try {
      await axios.post(`${API_URL}/mail/send`, { receiver, subject, body }, { headers: authHeaders });
      alert("Mail sent");
    } catch (err) {
      alert("Sending failed");
    }
  };

  return (
    <div>
      <h2>Send Mail</h2>
      <input type="email" placeholder="Receiver" value={receiver} onChange={(e) => setReceiver(e.target.value)} />
      <input type="text" placeholder="Subject" value={subject} onChange={(e) => setSubject(e.target.value)} />
      <textarea placeholder="Body" value={body} onChange={(e) => setBody(e.target.value)}></textarea>
      <button onClick={handleSend}>Send</button>
    </div>
  );
}

function Admin({ authHeaders }) {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    axios.get(`${API_URL}/admin/users`, { headers: authHeaders })
      .then((res) => setUsers(res.data))
      .catch(() => alert("Failed to load users"));
  }, [authHeaders]);

  const handleDelete = async (id) => {
    try {
      await axios.delete(`${API_URL}/admin/users/${id}`, { headers: authHeaders });
      setUsers(users.filter((user) => user.id !== id));
    } catch (err) {
      alert("Failed to delete user");
    }
  };

  return (
    <div>
      <h2>Admin Panel</h2>
      <ul>
        {users.map((user) => (
          <li key={user.id}>{user.email} <button onClick={() => handleDelete(user.id)}>Delete</button></li>
        ))}
      </ul>
    </div>
  );
}

export default App;
