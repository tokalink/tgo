document.getElementById('btn-fetch').addEventListener('click', async () => {
  const output = document.getElementById('output');
  output.textContent = 'Calling RPC...';
  try {
    const res = await fetch('/user.v1.UserService/GetProfile', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-Slug': 'demo_tenant',
        'Connect-Protocol-Version': '1'
      },
      body: JSON.stringify({ user_id: 'usr_frontend_demo' })
    });
    const data = await res.json();
    output.textContent = JSON.stringify(data, null, 2);
  } catch (err) {
    output.textContent = 'Error: ' + err.message;
  }
});