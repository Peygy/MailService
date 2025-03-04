import React, { useState, useEffect, useRef } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate, Navigate } from "react-router-dom";
import axios from "axios";
import { styled, createGlobalStyle } from 'styled-components';

const API_URL = "http://localhost:8081/api/v1";

const GlobalStyle = createGlobalStyle`
  body {
    background-color: #121212;
    color: #e0e0e0;
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 0;
  }

  a {
    color: #bb86fc;
    text-decoration: none;
    &:hover {
      color: #3700b3;
    }
  }

  button {
    background-color: #bb86fc;
    color: #121212;
    border: none;
    padding: 0.5rem 1rem;
    cursor: pointer;
    &:hover {
      background-color: #3700b3;
    }
  }

  input, textarea {
    background-color: #1e1e1e;
    color: #e0e0e0;
    border: 1px solid #bb86fc;
    padding: 0.5rem;
    font-size: 1rem;
    &:focus {
      outline: none;
      border-color: #3700b3;
    }
  }

  table {
    width: 100%;
    border-collapse: collapse;
    th, td {
      border: 1px solid #bb86fc;
      padding: 8px;
    }
    th {
      background-color: #1e1e1e;
    }
    tr:nth-child(even) {
      background-color: #1e1e1e;
    }
    tr:hover {
      background-color: #3700b3; 
    }
  }
`;

const NavBar = styled.nav`
  background-color: #1e1e1e;
  padding: 1rem;
  display: flex;
  justify-content: space-around;
  a {
    color: #bb86fc;
    text-decoration: none;
    font-size: 1.2rem;
    &:hover {
      color: #3700b3;
    }
  }
  button {
    background-color: #bb86fc;
    border: none;
    padding: 0.5rem 1rem;
    color: #121212;
    cursor: pointer;
    &:hover {
      background-color: #3700b3;
    }
  }
`;

const Container = styled.div`
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
`;

const Form = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
  input, textarea {
    background-color: #1e1e1e; 
    color: #e0e0e0; 
    border: 1px solid #bb86fc;
    padding: 0.5rem;
    font-size: 1rem;
    &:focus {
      outline: none;
      border-color: #3700b3; 
    }
  }
  button {
    background-color: #bb86fc; 
    color: #121212; 
    border: none;
    padding: 0.5rem;
    font-size: 1rem;
    cursor: pointer;
    &:hover {
      background-color: #3700b3; 
    }
  }
`;

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
  th, td {
    border: 1px solid #bb86fc; /* Фиолетовая рамка */
    padding: 8px;
  }
  th {
    background-color: #1e1e1e; /* Темный фон для заголовков */
  }
  tr:nth-child(even) {
    background-color: #1e1e1e; /* Темный фон для четных строк */
  }
  tr:hover {
    background-color: #3700b3; /* Темно-фиолетовый при наведении */
  }
`;

function App() {
  const [auth, setAuth] = useState(() => {
    const storedAuth = localStorage.getItem('auth');
    return storedAuth ? storedAuth : null;
  });

  const login = (username, password) => {
    const encoded = btoa(`${username}:${password}`);
    setAuth(encoded);
    localStorage.setItem('auth', encoded);
  };

  const logout = () => {
    setAuth(null);
    localStorage.removeItem('auth');
  };

  const authHeaders = auth ? { Authorization: `Basic ${auth}` } : {};

  return (
    <Router>
      <GlobalStyle />
      <NavBar>
        <Link to="/login">Login</Link>
        <Link to="/register">Register</Link>
        <Link to="/inbox">Inbox</Link>
        <Link to="/sent">Sent</Link>
        <Link to="/send">SendMail</Link>
        <Link to="/admin">Admin</Link>
        {auth && <button onClick={logout}>Logout</button>}
      </NavBar>
      <Container>
        <Routes>
          <Route path="/login" element={<Login onLogin={login} />} />
          <Route path="/register" element={<Register />} />
          <Route path="/inbox" element={auth ? <Inbox authHeaders={authHeaders} /> : <Navigate to="/login" />} />
          <Route path="/sent" element={auth ? <Sent authHeaders={authHeaders} /> : <Navigate to="/login" />} />
          <Route path="/send" element={auth ? <SendMail authHeaders={authHeaders} /> : <Navigate to="/login" />} />
          <Route path="/admin" element={auth ? <Admin authHeaders={authHeaders} /> : <Navigate to="/login" />} />
          <Route path="*" element={<Navigate to="/login" />} />
        </Routes>
      </Container>
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
      onLogin(email, password);
      navigate("/inbox");
    } catch (err) {
      alert("Login failed");
    }
  };

  return (
    <Form>
      <h2>Login</h2>
      <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button onClick={handleLogin}>Login</button>
    </Form>
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
    <Form>
      <h2>Register</h2>
      <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button onClick={handleRegister}>Register</button>
    </Form>
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
        <Table>
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
        </Table>
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
        <Table>
          <thead>
            <tr>
              <th>Receivers</th>
              <th>Subject</th>
              <th>Body</th>
              <th>Sent At</th>
            </tr>
          </thead>
          <tbody>
            {mails.map((mail) => (
              <tr key={mail.ID}>
                <td>{mail.Receivers}</td>
                <td>{mail.Subject}</td>
                <td>{mail.Body}</td>
                <td>{new Date(mail.CreatedAt).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </div>
  );
}

function SendMail({ authHeaders }) {
  const [receivers, setReceivers] = useState("");
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");

  const handleSend = async () => {
    const receiverList = receivers.split(',').map(email => email.trim());

    try {
      await axios.post(`${API_URL}/mail/send`, { receivers: receiverList, subject, body }, { headers: authHeaders });
      alert("Mail sent");
      setReceivers("");
      setSubject("");
      setBody("");
    } catch (err) {
      alert("Sending failed");
    }
  };

  return (
    <Form>
      <h2>Send Mail</h2>
      <input 
        type="text" 
        placeholder="Receivers (comma separated)" 
        value={receivers} 
        onChange={(e) => setReceivers(e.target.value)} 
      />
      <input 
        type="text" 
        placeholder="Subject" 
        value={subject} 
        onChange={(e) => setSubject(e.target.value)} 
      />
      <textarea 
        placeholder="Body" 
        value={body} 
        onChange={(e) => setBody(e.target.value)} 
      ></textarea>
      <button onClick={handleSend}>Send</button>
    </Form>
  );
}

function Admin({ authHeaders }) {
  const [users, setUsers] = useState([]);
  const isMounted = useRef(false);

  useEffect(() => {
    if (!isMounted.current && authHeaders.Authorization) {
      isMounted.current = true;
      axios.get(`${API_URL}/admin/users`, { headers: authHeaders })
        .then((res) => setUsers(res.data))
        .catch(() => alert("Failed to load users"));
    }
  }, [authHeaders]);

  const handleDelete = async (id) => {
    try {
      await axios.delete(`${API_URL}/admin/users/${id}`, { headers: authHeaders });
      setUsers(users.filter((user) => user.Id !== id));
    } catch (err) {
      alert("Failed to delete user");
    }
  };

  return (
    <div>
      <h2>Admin Panel</h2>
      <ul>
        {users.map((user) => (
          <li key={user.Id}>{user.Email} - {user.Role} <button onClick={() => handleDelete(user.Id)}>Delete</button></li>
        ))}
      </ul>
    </div>
  );
}

export default App;