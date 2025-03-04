const kill = require('kill-port');
const { exec } = require('child_process');

const PORT = 3000;

kill(PORT, 'tcp')
  .then(() => {
    console.log(`Порт ${PORT} освобожден. Запуск приложения...`);
    const process = exec('react-scripts start', { stdio: 'inherit' });

    process.stdout.on('data', (data) => {
      console.log(data);
    });

    process.stderr.on('data', (data) => {
      console.error(data);
    });
  })
  .catch((err) => {
    console.error('Не удалось освободить порт', err);
    process.exit(1);
  });