// --- Módulo de Audio Sintetizado nativo (Web Audio API) ---
class AudioManager {
  constructor() {
    this.ctx = null;
    this.noiseBuffer = null;
    this.volume = 0.5;
    this.muted = false;
    this.explosionType = 'both'; // 'both', 'synth', 'noise'
  }

  init() {
    if (!this.ctx) {
      this.ctx = new (window.AudioContext || window.webkitAudioContext)();
    }
    if (this.ctx.state === 'suspended') {
      this.ctx.resume();
    }
  }

  getGain(baseGain) {
    return this.muted ? 0 : baseGain * this.volume;
  }

  getNoiseBuffer() {
    if (this.noiseBuffer) return this.noiseBuffer;
    if (!this.ctx) return null;
    
    const bufferSize = this.ctx.sampleRate * 0.4;
    const buffer = this.ctx.createBuffer(1, bufferSize, this.ctx.sampleRate);
    const data = buffer.getChannelData(0);
    for (let i = 0; i < bufferSize; i++) {
      data[i] = Math.random() * 2 - 1;
    }
    this.noiseBuffer = buffer;
    return this.noiseBuffer;
  }

  playClick() {
    this.init();
    if (this.muted) return;
    const now = this.ctx.currentTime;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();

    osc.connect(gain);
    gain.connect(this.ctx.destination);

    osc.type = 'sine';
    osc.frequency.setValueAtTime(300, now);
    osc.frequency.exponentialRampToValueAtTime(600, now + 0.12);

    gain.gain.setValueAtTime(this.getGain(0.15), now);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.12);

    osc.start(now);
    osc.stop(now + 0.12);
  }

  playComplete() {
    this.init();
    if (this.muted) return;
    const now = this.ctx.currentTime;
    const notes = [523.25, 659.25, 783.99, 987.77];
    
    notes.forEach((freq, index) => {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.connect(gain);
      gain.connect(this.ctx.destination);

      osc.type = 'triangle';
      const triggerTime = now + (index * 0.12);

      osc.frequency.setValueAtTime(freq, triggerTime);
      gain.gain.setValueAtTime(0, triggerTime);
      gain.gain.linearRampToValueAtTime(this.getGain(0.2), triggerTime + 0.03);
      gain.gain.exponentialRampToValueAtTime(0.001, triggerTime + 0.6);

      osc.start(triggerTime);
      osc.stop(triggerTime + 0.65);
    });
  }

  playReveal() {
    this.init();
    if (this.muted) return;
    const now = this.ctx.currentTime;
    const notes = [349.23, 440.00, 523.25, 659.25, 783.99];

    notes.forEach((freq, index) => {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();

      osc.type = 'sine';
      osc.connect(gain);
      gain.connect(this.ctx.destination);

      const triggerTime = now + (index * 0.15);
      osc.frequency.setValueAtTime(freq, triggerTime);

      gain.gain.setValueAtTime(0, triggerTime);
      gain.gain.linearRampToValueAtTime(this.getGain(0.12), triggerTime + 0.08);
      gain.gain.exponentialRampToValueAtTime(0.001, triggerTime + 0.8);

      osc.start(triggerTime);
      osc.stop(triggerTime + 0.9);
    });
  }

  playExplosion(isAmbient = false) {
    this.init();
    if (this.muted) return;
    const now = this.ctx.currentTime;

    const useSynth = this.explosionType === 'both' || this.explosionType === 'synth';
    const useNoise = this.explosionType === 'both' || this.explosionType === 'noise';

    if (useSynth) {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      osc.connect(gain);
      gain.connect(this.ctx.destination);
      osc.type = 'sine';

      if (isAmbient) {
        osc.frequency.setValueAtTime(80, now);
        osc.frequency.exponentialRampToValueAtTime(10, now + 0.6);
        gain.gain.setValueAtTime(this.getGain(0.05), now);
        gain.gain.exponentialRampToValueAtTime(0.001, now + 0.6);
      } else {
        osc.frequency.setValueAtTime(140, now);
        osc.frequency.exponentialRampToValueAtTime(20, now + 0.45);
        gain.gain.setValueAtTime(this.getGain(0.22), now);
        gain.gain.exponentialRampToValueAtTime(0.001, now + 0.45);
      }
      osc.start(now);
      osc.stop(now + 0.65);
    }

    if (useNoise && !isAmbient) {
      try {
        const noise = this.ctx.createBufferSource();
        noise.buffer = this.getNoiseBuffer();

        const filter = this.ctx.createBiquadFilter();
        filter.type = 'bandpass';
        filter.frequency.value = 1600;
        filter.Q.value = 2.5;

        const noiseGain = this.ctx.createGain();
        noiseGain.gain.setValueAtTime(this.getGain(0.07), now);
        noiseGain.gain.exponentialRampToValueAtTime(0.001, now + 0.35);

        noise.connect(filter);
        filter.connect(noiseGain);
        noiseGain.connect(this.ctx.destination);

        noise.start(now);
        noise.stop(now + 0.4);
      } catch (e) {
        // No-op
      }
    }
  }
}

const audio = new AudioManager();

// --- Sincronización del Panel de Control de Sonido ---
const muteToggle = document.getElementById('muteToggle');
const volumeSlider = document.getElementById('volumeSlider');
const soundType = document.getElementById('soundType');
const soundConfig = document.getElementById('soundConfig');
const soundToggleBtn = document.getElementById('soundToggleBtn');

// Alternar visibilidad del panel de sonido
if (soundToggleBtn && soundConfig) {
  soundToggleBtn.addEventListener('click', (e) => {
    e.stopPropagation(); // Evita disparar fuegos artificiales al configurar
    soundConfig.classList.toggle('active');
  });
  
  // Opcional: Cerrar el panel al hacer clic fuera de él
  window.addEventListener('click', (e) => {
    if (soundConfig.classList.contains('active') && !soundConfig.contains(e.target) && e.target !== soundToggleBtn) {
      soundConfig.classList.remove('active');
    }
  });
}

if (muteToggle) {
  muteToggle.addEventListener('change', (e) => {
    audio.muted = !e.target.checked;
  });
  audio.muted = !muteToggle.checked;
}

if (volumeSlider) {
  volumeSlider.addEventListener('input', (e) => {
    audio.volume = parseFloat(e.target.value);
  });
  audio.volume = parseFloat(volumeSlider.value);
}

if (soundType) {
  soundType.addEventListener('change', (e) => {
    audio.explosionType = e.target.value;
  });
  audio.explosionType = soundType.value;
}

// Detener la propagación de eventos sobre el panel para no lanzar fuegos artificiales
if (soundConfig) {
  const stopPropagation = (e) => e.stopPropagation();
  soundConfig.addEventListener('pointerdown', stopPropagation);
  soundConfig.addEventListener('click', stopPropagation);
}

// --- Renderizado y Lógica del Canvas ---
const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');

function resizeCanvas() {
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;
}
window.addEventListener('resize', resizeCanvas);
resizeCanvas();

const gravity = 0.04;
const friction = 0.985;

let fireworks = [];
let particles = [];
let targetReached = false;
let clickCount = 0;
const clicksRequired = 5;

class Firework {
  constructor(sx, sy, tx, ty, customColor = null, isManual = false) {
    this.x = sx;
    this.y = sy;
    this.tx = tx;
    this.ty = ty;
    this.distanceToTarget = calculateDistance(sx, sy, tx, ty);
    this.distanceTraveled = 0;
    this.coordinates = [];
    this.coordinateCount = 3;
    while(this.coordinateCount--) {
      this.coordinates.push([this.x, this.y]);
    }
    this.angle = Math.atan2(ty - sy, tx - sx);
    this.speed = 2.5;
    this.acceleration = 1.04;
    this.brightness = Math.random() * 20 + 60;
    this.customColor = customColor;
    this.isManual = isManual;
  }

  update(index) {
    this.coordinates.pop();
    this.coordinates.unshift([this.x, this.y]);
    
    this.speed *= this.acceleration;
    
    let vx = Math.cos(this.angle) * this.speed;
    let vy = Math.sin(this.angle) * this.speed;
    this.distanceTraveled = calculateDistance(this.x, this.y, this.x + vx, this.y + vy);
    
    if(this.distanceTraveled >= this.distanceToTarget) {
      createParticles(this.tx, this.ty, this.customColor);
      
      if (this.isManual) {
        fireworks.push(new Firework(Math.random() * canvas.width, canvas.height, this.tx, this.ty, 'rgba(255, 42, 116, 1)', false));
        audio.playExplosion(false);
      } else {
        if (Math.random() < 0.35) {
          audio.playExplosion(true);
        }
      }

      fireworks.splice(index, 1);
    } else {
      this.x += vx;
      this.y += vy;
    }
  }

  draw() {
    ctx.beginPath();
    ctx.moveTo(this.coordinates[this.coordinates.length - 1][0], this.coordinates[this.coordinates.length - 1][1]);
    ctx.lineTo(this.x, this.y);
    ctx.strokeStyle = this.customColor ? this.customColor : `hsl(${Math.random() * 360}, 100%, ${this.brightness}%)`;
    ctx.lineWidth = 1.8;
    ctx.stroke();
  }
}

class Particle {
  constructor(x, y, customColor = null) {
    this.x = x;
    this.y = y;
    this.coordinates = [];
    this.coordinateCount = 5;
    while(this.coordinateCount--) {
      this.coordinates.push([this.x, this.y]);
    }
    this.angle = Math.random() * Math.PI * 2;
    this.speed = Math.random() * 8 + 1;
    this.friction = friction;
    this.gravity = gravity;
    this.hue = Math.random() * 360;
    this.customColor = customColor;
    this.decay = Math.random() * 0.015 + 0.012;
    this.alpha = 1;
  }

  update(index) {
    this.coordinates.pop();
    this.coordinates.unshift([this.x, this.y]);
    this.speed *= this.friction;
    this.x += Math.cos(this.angle) * this.speed;
    this.y += Math.sin(this.angle) * this.speed + this.gravity;
    this.alpha -= this.decay;
    
    if(this.alpha <= this.decay) {
      particles.splice(index, 1);
    }
  }

  draw() {
    ctx.beginPath();
    ctx.moveTo(this.coordinates[this.coordinates.length - 1][0], this.coordinates[this.coordinates.length - 1][1]);
    ctx.lineTo(this.x, this.y);
    if (this.customColor) {
      ctx.strokeStyle = this.customColor.replace('1)', `${this.alpha})`);
    } else {
      ctx.strokeStyle = `hsla(${this.hue}, 100%, 60%, ${this.alpha})`;
    }
    ctx.lineWidth = 1.5;
    ctx.stroke();
  }
}

function calculateDistance(x1, y1, x2, y2) {
  return Math.sqrt(Math.pow(x1 - x2, 2) + Math.pow(y1 - y2, 2));
}

function createParticles(x, y, customColor = null) {
  let particleCount = customColor ? 150 : 75;
  while(particleCount--) {
    particles.push(new Particle(x, y, customColor));
  }
}

const giftButton = document.getElementById('giftButton');
const giftWrapper = document.getElementById('giftWrapper');
const giftHint = document.getElementById('giftHint');
const congratsInterface = document.getElementById('congratsInterface');
const backBtn = document.getElementById('backBtn');

const fillStop = document.getElementById('fillStop');
const emptyStop = document.getElementById('emptyStop');

if (backBtn) {
  backBtn.addEventListener('pointerdown', (e) => {
    e.stopPropagation();
  });
}

giftButton.addEventListener('click', (e) => {
  e.stopPropagation();
  if (targetReached) return;

  clickCount++;
  audio.playClick();

  const progressRatio = clickCount / clicksRequired;
  const fillPercentage = progressRatio * 100;

  if (fillStop && emptyStop) {
    fillStop.setAttribute('offset', `${fillPercentage}%`);
    emptyStop.setAttribute('offset', `${fillPercentage}%`);
  }

  const rect = giftButton.getBoundingClientRect();
  const clickX = rect.left + rect.width / 2;
  const clickY = rect.top + rect.height / 2;

  giftButton.style.transform = `scale(${1 + progressRatio * 0.25})`;

  fireworks.push(new Firework(Math.random() * canvas.width, canvas.height, clickX, clickY, 'rgba(255, 42, 116, 1)'));

  if (clickCount >= clicksRequired) {
    targetReached = true;
    audio.playComplete();

    giftButton.classList.add('happy');
    giftHint.innerText = "¡Completado! ❤️";

    setTimeout(() => {
      createParticles(clickX, clickY, 'rgba(255, 42, 116, 1)');
      createParticles(clickX, clickY, 'rgba(0, 255, 255, 1)');
      createParticles(clickX, clickY, 'rgba(251, 243, 140, 1)');
      
      giftWrapper.style.transition = 'opacity 1s ease';
      giftWrapper.style.opacity = '0';
      
      setTimeout(() => {
        giftWrapper.style.display = 'none';
        congratsInterface.classList.add('show');
        congratsInterface.style.pointerEvents = 'auto';

        audio.playReveal();

        setTimeout(() => {
          if (backBtn) {
            backBtn.style.opacity = '1';
            backBtn.style.pointerEvents = 'auto';
          }
        }, 3000);

      }, 1000);
    }, 600);
  } else {
    giftHint.innerText = "¡Sigue presionando!";
  }
});

window.addEventListener('pointerdown', (e) => {
  if (!targetReached) return;
  
  const x = e.clientX;
  const y = e.clientY;
  fireworks.push(new Firework(x, canvas.height, x, y, null, true));
});

let autoLaunchTimer = 0;
function autoLaunch() {
  autoLaunchTimer++;
  const interval = targetReached ? 25 : 80;
  
  if (autoLaunchTimer % interval === 0) {
    const startX = Math.random() * canvas.width;
    const targetX = Math.random() * canvas.width;
    const targetY = Math.random() * (canvas.height * 0.6);
    
    if (targetReached) {
      fireworks.push(new Firework(startX, canvas.height, targetX, targetY, null, false));
    } else if (autoLaunchTimer % 160 === 0) {
      fireworks.push(new Firework(startX, canvas.height, targetX, targetY * 1.2, null, false));
    }
  }
}

function loop() {
  requestAnimationFrame(loop);
  
  ctx.globalCompositeOperation = 'destination-out';
  ctx.fillStyle = 'rgba(5, 5, 16, 0.25)';
  ctx.fillRect(0, 0, canvas.width, canvas.height);
  
  ctx.globalCompositeOperation = 'lighter';
  
  let i = fireworks.length;
  while(i--) {
    fireworks[i].draw();
    fireworks[i].update(i);
  }
  
  let j = particles.length;
  while(j--) {
    particles[j].draw();
    particles[j].update(j);
  }

  autoLaunch();
}

loop();