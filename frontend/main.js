// Application State
const state = {
  token: localStorage.getItem('tgo_jwt') || null,
  user: JSON.parse(localStorage.getItem('tgo_user') || 'null'),
  activeTenant: localStorage.getItem('tgo_tenant') || 'acme_corp',
  currentTab: 'overview'
};

// DOM Elements
const loginView = document.getElementById('login-view');
const dashboardView = document.getElementById('dashboard-view');
const loginForm = document.getElementById('login-form');
const emailInput = document.getElementById('email-input');
const passwordInput = document.getElementById('password-input');
const tenantSelect = document.getElementById('tenant-select');
const btnLogout = document.getElementById('btn-logout');
const activeTenantSwitch = document.getElementById('active-tenant-switch');
const userNameDisplay = document.getElementById('user-name-display');
const userRoleDisplay = document.getElementById('user-role-display');
const userAvatar = document.getElementById('user-avatar');
const breadcrumbCurrent = document.getElementById('breadcrumb-current');
const topLatency = document.getElementById('top-latency');

// Init Lifecycle
document.addEventListener('DOMContentLoaded', () => {
  setupEventListeners();
  if (state.token && state.user) {
    showDashboard();
  } else {
    showLogin();
  }
});

function setupEventListeners() {
  // Quick Fill Demo Credentials
  document.querySelectorAll('.quick-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      const tenant = btn.getAttribute('data-tenant');
      const email = btn.getAttribute('data-email');
      tenantSelect.value = tenant;
      emailInput.value = email;
      passwordInput.value = 'tgo-demo-password';
    });
  });

  // Login Submit
  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
    const tenant_slug = tenantSelect.value;

    const btn = document.getElementById('btn-login');
    btn.querySelector('.btn-text').textContent = 'Authenticating...';

    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, tenant_slug })
      });

      if (!res.ok) {
        throw new Error('Login failed: ' + res.statusText);
      }

      const data = await res.json();
      state.token = data.token;
      state.user = data.user;
      state.activeTenant = data.user.tenant_slug || tenant_slug;

      localStorage.setItem('tgo_jwt', state.token);
      localStorage.setItem('tgo_user', JSON.stringify(state.user));
      localStorage.setItem('tgo_tenant', state.activeTenant);

      showDashboard();
    } catch (err) {
      alert(err.message || 'Authentication error');
    } finally {
      btn.querySelector('.btn-text').textContent = 'Sign In to Dashboard';
    }
  });

  // Logout
  btnLogout.addEventListener('click', () => {
    state.token = null;
    state.user = null;
    localStorage.removeItem('tgo_jwt');
    localStorage.removeItem('tgo_user');
    showLogin();
  });

  // Navigation Tabs
  document.querySelectorAll('.sidebar-nav .nav-item').forEach(btn => {
    btn.addEventListener('click', () => {
      const tab = btn.getAttribute('data-tab');
      switchTab(tab);
    });
  });

  // Active Tenant Switcher
  activeTenantSwitch.addEventListener('change', (e) => {
    state.activeTenant = e.target.value;
    localStorage.setItem('tgo_tenant', state.activeTenant);
    fetchDashboardStats();
  });

  // ConnectRPC Tester Button
  const btnCallRpc = document.getElementById('btn-call-rpc');
  if (btnCallRpc) {
    btnCallRpc.addEventListener('click', callConnectRPC);
  }

  // Refresh Stats Button
  const btnRefreshStats = document.getElementById('btn-refresh-stats');
  if (btnRefreshStats) {
    btnRefreshStats.addEventListener('click', fetchDashboardStats);
  }
}

function showLogin() {
  loginView.classList.remove('hidden');
  dashboardView.classList.add('hidden');
}

function showDashboard() {
  loginView.classList.add('hidden');
  dashboardView.classList.remove('hidden');

  // Populate user metadata
  if (state.user) {
    userNameDisplay.textContent = state.user.name || 'Admin User';
    userRoleDisplay.textContent = `${state.user.role || 'Admin'} • ${state.activeTenant}`;
    const initials = (state.user.name || 'Admin').split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
    userAvatar.textContent = initials || 'TG';
  }

  activeTenantSwitch.value = state.activeTenant;
  fetchDashboardStats();
}

function switchTab(tabId) {
  state.currentTab = tabId;
  document.querySelectorAll('.sidebar-nav .nav-item').forEach(btn => {
    btn.classList.toggle('active', btn.getAttribute('data-tab') === tabId);
  });

  document.querySelectorAll('.tab-content').forEach(tab => {
    tab.classList.remove('active');
  });

  const activeContent = document.getElementById(`tab-${tabId}`);
  if (activeContent) {
    activeContent.classList.add('active');
  }

  const tabLabels = {
    overview: 'Dashboard Overview',
    tenants: 'Tenants & Schema Isolation',
    rpc: 'ConnectRPC Studio',
    team: 'Team Directory',
    database: 'Schema Migrations',
    settings: 'Tenant Settings'
  };

  breadcrumbCurrent.textContent = tabLabels[tabId] || 'Dashboard';
}

async function fetchDashboardStats() {
  const startTime = performance.now();
  try {
    const res = await fetch('/api/dashboard/stats', {
      headers: {
        'Authorization': `Bearer ${state.token}`,
        'X-Tenant-Slug': state.activeTenant
      }
    });

    const duration = (performance.now() - startTime).toFixed(1);
    topLatency.textContent = `${duration} ms (ConnectRPC & DB)`;

    if (!res.ok) return;
    const data = await res.json();

    // Populate Overview Stats
    document.getElementById('stat-revenue').textContent = data.total_revenue || '$128,450.00';
    document.getElementById('stat-users').textContent = (data.active_users || 2840).toLocaleString();
    document.getElementById('stat-throughput').textContent = data.rpc_throughput || '172,500 req/s';
    document.getElementById('stat-schema').textContent = data.schema_name || `tenant_${state.activeTenant}`;

    // Populate Activity Table
    const eventsTbody = document.getElementById('events-tbody');
    if (eventsTbody && data.recent_activity) {
      eventsTbody.innerHTML = data.recent_activity.map(evt => `
        <tr>
          <td><code style="color: #38bdf8">${evt.id}</code></td>
          <td><strong>${evt.event}</strong></td>
          <td><span style="color: #94a3b8">${evt.tenant}</span></td>
          <td><span style="color: #64748b">${evt.time}</span></td>
          <td><span class="badge-status success">${evt.status}</span></td>
        </tr>
      `).join('');
    }

    // Populate Team Table
    const teamTbody = document.getElementById('team-tbody');
    if (teamTbody && data.team_members) {
      teamTbody.innerHTML = data.team_members.map(member => `
        <tr>
          <td><strong>${member.name}</strong></td>
          <td><span style="color: #94a3b8">${member.email}</span></td>
          <td><span class="badge-status success">${member.role}</span></td>
          <td><span style="color: #10b981">● ${member.status}</span></td>
        </tr>
      `).join('');
    }
  } catch (err) {
    console.error('Failed to fetch dashboard stats:', err);
  }
}

async function callConnectRPC() {
  const userIdInput = document.getElementById('rpc-user-id');
  const rpcResult = document.getElementById('rpc-result');
  const rpcTime = document.getElementById('rpc-time');

  const userId = userIdInput.value.trim() || 'usr_frontend_demo';
  rpcResult.textContent = 'Calling ConnectRPC service...';
  rpcTime.textContent = 'Calling...';

  const startTime = performance.now();
  try {
    const res = await fetch('/user.v1.UserService/GetProfile', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Connect-Protocol-Version': '1',
        'X-Tenant-Slug': state.activeTenant,
        'Authorization': `Bearer ${state.token}`
      },
      body: JSON.stringify({ user_id: userId })
    });

    const elapsed = (performance.now() - startTime).toFixed(2);
    rpcTime.textContent = `${elapsed} ms`;

    const data = await res.json();
    rpcResult.textContent = JSON.stringify(data, null, 2);
  } catch (err) {
    rpcResult.textContent = JSON.stringify({
      error: err.message,
      status: 'RPC_FAILED'
    }, null, 2);
    rpcTime.textContent = 'Error';
  }
}