import React, { useState, useEffect, useRef } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate, Navigate } from "react-router-dom";
import axios from "axios";
import { styled, createGlobalStyle, keyframes } from 'styled-components';

const API_URL = "/api/v1";

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

const fadeIn = keyframes`
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
`;

const fadeOut = keyframes`
  from { opacity: 1; transform: scale(1); }
  to { opacity: 0; transform: scale(0.95); }
`;

const MailPopup = styled.div`
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) !important;;
  background-color: #1e1e1e;
  border: 2px solid #bb86fc;
  padding: 20px;
  border-radius: 8px;
  z-index: 1000;
  max-width: 600px;
  width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
  animation: ${props => props.$isClosing ? fadeOut : fadeIn} 0.3s ease-in-out;
  
  h3 {
    margin-top: 0;
    color: #bb86fc;
    border-bottom: 1px solid #bb86fc;
    padding-bottom: 10px;
  }
  
  p {
    white-space: pre-wrap;
    word-break: break-word;
  }
  
  button {
    position: absolute;
    top: 10px;
    right: 10px;
    background: none;
    border: none;
    color: #bb86fc;
    font-size: 20px;
    cursor: pointer;
  }
`;

const Overlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  z-index: 999;
`;

const NotificationContainer = styled.div`
  position: fixed;
  top: 20px;
  right: 20px;
  background-color: #ff4444;
  color: white;
  padding: 10px 20px;
  border-radius: 5px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 1000;
`;

const CloseButton = styled.button`
  background: none;
  border: none;
  color: white;
  font-size: 16px;
  cursor: pointer;
`;

const Notification = ({ message, onClose }) => {
  useEffect(() => {
    const timer = setTimeout(() => {
      onClose();
    }, 2000);
    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <NotificationContainer>
      <span>{message}</span>
      <CloseButton onClick={onClose}>×</CloseButton>
    </NotificationContainer>
  );
};

const Header = styled.header`
  background-color: #1e1e1e;
  padding: 1rem;
  display: flex;
  justify-content: flex-end;
  align-items: center;
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

const Layout = styled.div`
  display: flex;
  height: calc(100vh - 60px);
`;

const SideNav = styled.nav`
  background-color: #1e1e1e;
  width: 200px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  a {
    color: #bb86fc;
    text-decoration: none;
    font-size: 1.2rem;
    &:hover {
      color: #3700b3;
    }
  }
`;

const MainContent = styled.div`
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
`;

const Container = styled.div`
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
    cursor: pointer;
  }
  .sortable {
    cursor: pointer;
  }

  .sortable:hover {
    text-decoration: underline;
  }
`;

function App() {
  const [auth, setAuth] = useState(() => {
    const storedAuth = localStorage.getItem('auth');
    return storedAuth ? storedAuth : null;
  });
  const [notification, setNotification] = useState(null);

  const showNotification = (message) => {
    setNotification(message);
  };

  const closeNotification = () => {
    setNotification(null);
  };

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

  const getEmailFromAuth = () => {
    if (!auth) return null;
    const decoded = atob(auth);
    return decoded.split(':')[0];
  };

  const email = getEmailFromAuth();

  const isAdmin = email && email.includes('@admin.gomail.kurs');

  return (
    <Router>
      <GlobalStyle />
      {notification && <Notification message={notification} onClose={closeNotification} />}
      {auth ? (
        <>
          <Header>
            <button onClick={logout}>Выйти</button>
          </Header>
          <Layout>
            <SideNav>
              <Link to="/inbox">Входящие</Link>
              <Link to="/sent">Отправленные</Link>
              <Link to="/trash">Корзина</Link>
              <Link to="/send">Отправить письмо</Link>
              {isAdmin && (
                    <>
                    <h3>Админка</h3>
                    <Link to="/admin/mails">Письма</Link>
                    <Link to="/admin/users">Пользователи</Link>
                  </>
              )}
            </SideNav>
            <MainContent>
              <Container>
                <Routes>
                  <Route path="/inbox" element={<Inbox authHeaders={authHeaders} showNotification={showNotification} />} />
                  <Route path="/sent" element={<Sent authHeaders={authHeaders} showNotification={showNotification} />} />
                  <Route path="/trash" element={<Trash authHeaders={authHeaders} showNotification={showNotification} />} />
                  <Route path="/send" element={<SendMail authHeaders={authHeaders} showNotification={showNotification} />} />
                  {isAdmin && (
                    <>
                      <Route path="/admin/mails" element={<MailAdmin authHeaders={authHeaders} showNotification={showNotification} />} />
                      <Route path="/admin/users" element={<UserAdmin authHeaders={authHeaders} showNotification={showNotification} />} />
                    </>
                  )}
                  <Route path="*" element={<Navigate to="/inbox" />} />
                </Routes>
              </Container>
            </MainContent>
          </Layout>
        </>
      ) : (
        <Container>
          <Routes>
            <Route path="/login" element={<Login onLogin={login} showNotification={showNotification} />} />
            <Route path="/register" element={<Register onReg={login} showNotification={showNotification} />} />
            <Route path="*" element={<Navigate to="/login" />} />
          </Routes>
        </Container>
      )}
    </Router>
  );
}

function Login({ onLogin, showNotification }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      await axios.post(`${API_URL}/login`, { email, password });
      onLogin(email, password);
      navigate("/inbox");
    } catch (err) {
      showNotification(err.response?.data?.message || "Ошибка входа");
    }
  };

  return (
    <Form>
      <h2>Вход</h2>
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <div style={{ position: "relative", width: "100%" }}>
        <input
          type={showPassword ? "text" : "password"}
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          style={{ width: "98%" }}
        />
        <button
          onClick={() => setShowPassword(!showPassword)}
          style={{
            position: "absolute",
            right: "10px",
            top: "50%",
            transform: "translateY(-50%)",
            background: "none",
            border: "none",
            cursor: "pointer",
            color: "#bb86fc",
            padding: "0",
          }}
        >
          {showPassword ? "👁️" : "👁️‍🗨️"}
        </button>
      </div>
      <button onClick={handleLogin}>Войти</button>
      <p>
        Нет аккаунта? <Link to="/register">Регистрация</Link>
      </p>
    </Form>
  );
}

function Register({ onReg, showNotification }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();

  const validateEmail = (email) => {
    const regex = /^[^@]+@gomail\.kurs$/;
    return regex.test(email);
  };

  const handleRegister = async () => {
    setError("");

    if (!validateEmail(email)) {
      setError("Email должен содержать один @ и заканчиваться на gomail.kurs");
      return;
    }

    if (password !== confirmPassword) {
      setError("Пароли не совпадают");
      return;
    }

    try {
      await axios.post(`${API_URL}/register`, { email, password });
      onReg(email, password);
      navigate("/inbox");
    } catch (err) {
      showNotification(err.response?.data?.message || "Ошибка регистрации");
    }
  };

  return (
    <Form>
      <h2>Регистрация</h2>
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => {
          setEmail(e.target.value);
          setError("");
        }}
      />
      {error && <p style={{ color: "red" }}>{error}</p>}
      <div style={{ position: "relative" }}>
        <input
          type={showPassword ? "text" : "password"}
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          style={{ width: "98%" }}
        />
        <button
          onClick={() => setShowPassword(!showPassword)}
          style={{
            position: "absolute",
            right: "10px",
            top: "50%",
            transform: "translateY(-50%)",
            background: "none",
            border: "none",
            cursor: "pointer",
            color: "#bb86fc",
          }}
        >
          {showPassword ? "👁️" : "👁️‍🗨️"}
        </button>
      </div>
      <input
        type={showPassword ? "text" : "password"}
        placeholder="Подтвердите пароль"
        value={confirmPassword}
        onChange={(e) => setConfirmPassword(e.target.value)}
      />
      <button onClick={handleRegister}>Регистрация</button>
      <p>
        Уже есть аккаунт? <Link to="/login">Войти</Link>
      </p>
    </Form>
  );
}

function MailDetailsPopup({ mail, onClose }) {
  const [isClosing, setIsClosing] = useState(false);

  const handleClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      onClose();
    }, 300);
  };

  return (
    <>
      <Overlay onClick={handleClose} />
      <MailPopup $isClosing={isClosing}>
        <button onClick={handleClose}>×</button>
        <h3>{mail.Subject}</h3>
        <p>{mail.Body}</p>
      </MailPopup>
    </>
  );
}

function Inbox({ authHeaders, showNotification }) {
  const [mails, setMails] = useState([]);
  const [selectedMail, setSelectedMail] = useState(null);
  const [lastClickTime, setLastClickTime] = useState(0);

  useEffect(() => {
    axios.get(`${API_URL}/mail/inbox`, { headers: authHeaders })
      .then((res) => setMails(res.data.mails))
      .catch(() => showNotification("Ошибка загрузки входящих"));
  }, [authHeaders, showNotification]);

  const handleArchive = async (mailID) => {
    try {
      await axios.post(`${API_URL}/mail/${mailID}/archive`, {}, { headers: authHeaders });
      setMails(mails.filter((mail) => mail.ID !== mailID));
    } catch {
      showNotification("Ошибка архивации письма");
    }
  };

  const extractEmail = (sender) => {
    const match = sender.match(/<([^>]+)>$/);
    return match ? match[1] : sender;
  };

  const handleRowClick = (mail, e) => {
    const now = new Date().getTime();
    if (now - lastClickTime < 300) {
      setSelectedMail(mail);
    }
    setLastClickTime(now);
  };

  return (
    <div>
      <h2>Входящие</h2>
      {mails.length === 0 ? (
        <p>Нет писем</p>
      ) : (
        <Table>
          <thead>
            <tr>
              <th>Отправитель</th>
              <th>Дата получения</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {mails.map((mail) => (
              <tr key={mail.ID} onClick={(e) => handleRowClick(mail, e)}>
                <td>{extractEmail(mail.Sender)}</td>
                <td>{new Date(mail.CreatedAt).toLocaleString()}</td>
                <td>
                  <button onClick={(e) => {
                    e.stopPropagation();
                    handleArchive(mail.ID);
                  }}>В корзину</button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
      {selectedMail && (
        <MailDetailsPopup 
          mail={selectedMail} 
          onClose={() => setSelectedMail(null)} 
        />
      )}
    </div>
  );
};

function Sent({ authHeaders, showNotification }) {
  const [mails, setMails] = useState([]);
  const [selectedMail, setSelectedMail] = useState(null);
  const [lastClickTime, setLastClickTime] = useState(0);

  useEffect(() => {
    axios.get(`${API_URL}/mail/sent`, { headers: authHeaders })
      .then((res) => setMails(res.data.mails))
      .catch(() => showNotification("Ошибка загрузки отправленных"));
  }, [authHeaders, showNotification]);

  const handleRowClick = (mail, e) => {
    const now = new Date().getTime();
    if (now - lastClickTime < 300) {
      setSelectedMail(mail);
    }
    setLastClickTime(now);
  };

  return (
    <div>
      <h2>Отправленные</h2>
      {mails.length === 0 ? (
        <p>Нет отправленных писем</p>
      ) : (
        <Table>
          <thead>
            <tr>
              <th>Получатели</th>
              <th>Дата отправки</th>
            </tr>
          </thead>
          <tbody>
            {mails.map((mail) => (
              <tr key={mail.ID} onClick={(e) => handleRowClick(mail, e)}>
                <td>{mail.Receivers}</td>
                <td>{new Date(mail.CreatedAt).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
      {selectedMail && (
        <MailDetailsPopup 
          mail={selectedMail} 
          onClose={() => setSelectedMail(null)} 
        />
      )}
    </div>
  );
}

function SendMail({ authHeaders, showNotification }) {
  const [receivers, setReceivers] = useState("");
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");

  const handleSend = async () => {
    const receiverList = receivers.split(',').map(email => email.trim());

    try {
      await axios.post(`${API_URL}/mail/send`, { receivers: receiverList, subject, body }, { headers: authHeaders });
      setReceivers("");
      setSubject("");
      setBody("");
    } catch (err) {
      showNotification(err.response?.data?.message || "Ошибка отправки");
    }
  };

  return (
    <Form>
      <h2>Отправить письмо</h2>
      <input 
        type="text" 
        placeholder="Получатели (через запятую)" 
        value={receivers} 
        onChange={(e) => setReceivers(e.target.value)} 
      />
      <input 
        type="text" 
        placeholder="Тема" 
        value={subject} 
        onChange={(e) => setSubject(e.target.value)} 
      />
      <textarea 
        placeholder="Содержание" 
        value={body} 
        onChange={(e) => setBody(e.target.value)} 
      ></textarea>
      <button onClick={handleSend}>Отправить</button>
    </Form>
  );
}

function UserAdmin({ authHeaders, showNotification }) {
  const [users, setUsers] = useState([]);
  const isMounted = useRef(false);

  useEffect(() => {
    if (!isMounted.current && authHeaders.Authorization) {
      isMounted.current = true;
      axios.get(`${API_URL}/admin/users`, { headers: authHeaders })
        .then((res) => setUsers(res.data))
        .catch(() => showNotification("Ошибка загрузки пользователей"));
    }
  }, [authHeaders, showNotification]);

  const handleDelete = async (id) => {
    try {
      await axios.delete(`${API_URL}/admin/users/${id}`, { headers: authHeaders });
      setUsers(users.filter((user) => user.Id !== id));
    } catch (err) {
      showNotification(err.response?.data?.message || "Ошибка удаления пользователя");
    }
  };

  return (
    <div>
      <h2>Управление пользователями</h2>
      <ul>
        {users.map((user) => (
          <li key={user.Id}>{user.Email} - {user.Role} <button onClick={() => handleDelete(user.Id)}>Удалить</button></li>
        ))}
      </ul>
    </div>
  );
}

function MailAdmin({ authHeaders, showNotification }) {
  const [mails, setMails] = useState([]);
  const [filter, setFilter] = useState("");
  const [sortField, setSortField] = useState("CreatedAt");
  const [sortOrder, setSortOrder] = useState("desc");
  const [selectedMail, setSelectedMail] = useState(null);
  const [lastClickTime, setLastClickTime] = useState(0);

  useEffect(() => {
    if (authHeaders.Authorization) {
      axios
        .get(`${API_URL}/admin/mails`, { headers: authHeaders })
        .then((res) => {
          setMails(res.data || []);
        })
        .catch(() => showNotification("Ошибка загрузки писем"));
    }
  }, [authHeaders]);

  const handleDelete = async (id) => {
    try {
      await axios.delete(`${API_URL}/admin/mails/${id}`, { headers: authHeaders });
      setMails(mails.filter((mail) => mail.ID !== id));
    } catch (err) {
      showNotification(err.response?.data?.message || "Ошибка удаления письма");
    }
  };

  const filteredMails = mails.filter((mail) =>
    [mail.Sender, mail.Subject, mail.Body].some((field) =>
      field.toLowerCase().includes(filter.toLowerCase())
    )
  );

  const sortedMails = [...filteredMails].sort((a, b) => {
    if (sortField === "CreatedAt") {
      return sortOrder === "asc"
        ? new Date(a[sortField]) - new Date(b[sortField])
        : new Date(b[sortField]) - new Date(a[sortField]);
    }
    return sortOrder === "asc"
      ? a[sortField].localeCompare(b[sortField])
      : b[sortField].localeCompare(a[sortField]);
  });

  const handleSort = (field) => {
    if (sortField === field) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      setSortOrder("asc");
    }
  };

  const extractEmail = (sender) => {
    const match = sender.match(/<([^>]+)>$/);
    return match ? match[1] : sender;
  };

  const handleRowClick = (mail, e) => {
    const now = new Date().getTime();
    if (now - lastClickTime < 300) {
      setSelectedMail(mail);
    }
    setLastClickTime(now);
  };

  return (
    <div>
      <h2>Управление письмами</h2>
      <input 
        type="text"
        placeholder="Фильтр"
        value={filter}
        onChange={(e) => setFilter(e.target.value)}
      />
      {mails.length === 0 ? (
        <p>Нет писем</p>
      ) : (
        <Table>
          <thead>
            <tr>
              <th
                className="sortable"
                onClick={() => handleSort("Sender")}
              >
                Отправитель
              </th>
              <th>Дата</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {sortedMails.map((mail) => (
              <tr key={mail.ID} onClick={(e) => handleRowClick(mail, e)}>
                <td>{extractEmail(mail.Sender)}</td>
                <td>{new Date(mail.CreatedAt).toLocaleString()}</td>
                <td>
                  <button onClick={(e) => {
                    e.stopPropagation();
                    handleDelete(mail.ID);
                  }}>Удалить</button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
      {selectedMail && (
        <MailDetailsPopup 
          mail={selectedMail} 
          onClose={() => setSelectedMail(null)} 
        />
      )}
    </div>
  );
}

const Trash = ({ authHeaders, showNotification }) => {
  const [mails, setMails] = useState([]);
  const [selectedMail, setSelectedMail] = useState(null);
  const [lastClickTime, setLastClickTime] = useState(0);

  const fetchTrash = () => {
    axios.post(`${API_URL}/mail/trash`, {}, { headers: authHeaders })
      .then((res) => setMails(res.data.mails || []))
      .catch(() => showNotification("Ошибка загрузки корзины"));
  };

  useEffect(() => {
    fetchTrash();
  }, []);

  const handleRestore = async (mailID) => {
    try {
      await axios.post(`${API_URL}/mail/${mailID}/unarchive`, {}, { headers: authHeaders });
      setMails(mails.filter((mail) => mail.ID !== mailID));;
    } catch {
      showNotification("Ошибка восстановления письма");
    }
  };

  const handleDelete = async (mailID) => {
    try {
      await axios.delete(`${API_URL}/mail/${mailID}/delete`, { headers: authHeaders });
      fetchTrash();
    } catch {
      showNotification("Ошибка удаления письма");
    }
  };

  const handleDeleteAll = async () => {
    try {
      await Promise.all(mails.map(mail => 
        axios.delete(`${API_URL}/mail/${mail.ID}/delete`, { headers: authHeaders })
      ));
      setMails([]);
    } catch {
      showNotification("Ошибка очистки корзины");
    }
  };

  const handleRowClick = (mail, e) => {
    const now = new Date().getTime();
    if (now - lastClickTime < 300) {
      setSelectedMail(mail);
    }
    setLastClickTime(now);
  };

  return (
    <div>
      <h2>Корзина</h2>
      {mails.length === 0 ? (
        <p>Корзина пуста</p>
      ) : (
        <>
          <button onClick={handleDeleteAll}>Очистить корзину</button>
          <Table>
            <thead>
              <tr>
                <th>Тема</th>
                <th>Дата</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {mails.map((mail) => (
                <tr key={mail.ID} onClick={(e) => handleRowClick(mail, e)}>
                  <td>{mail.Subject}</td>
                  <td>{new Date(mail.CreatedAt).toLocaleString()}</td>
                  <td>
                    <button onClick={(e) => {
                      e.stopPropagation();
                      handleRestore(mail.ID);
                    }}>Восстановить</button>
                    <button onClick={(e) => {
                      e.stopPropagation();
                      handleDelete(mail.ID);
                    }} style={{ marginLeft: "8px", backgroundColor: "red" }}>
                      Удалить
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </>
      )}
      {selectedMail && (
        <MailDetailsPopup 
          mail={selectedMail} 
          onClose={() => setSelectedMail(null)} 
        />
      )}
    </div>
  );
};

export default App;