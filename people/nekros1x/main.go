package main

import (
	"log"
	"net/http"
	"os"
)

const page = `<!doctype html>
<html lang="es">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">
<title>Feliz cumple 三郎</title>
<link href="https://fonts.googleapis.com/css2?family=Baloo+2:wght@500;700;800&display=swap" rel="stylesheet">
<style>
  :root{
    --noche:#1a1333; --noche2:#2a1a4a;
    --frutilla:#ff5d8f; --crema:#ffe8c9;
    --dorado:#ffc857; --llama:#ff9f1c; --menta:#7bdff2;
  }
  *{margin:0;padding:0;box-sizing:border-box}
  html,body{height:100%;overflow:hidden;font-family:'Baloo 2',system-ui,sans-serif;
    background:radial-gradient(ellipse at 50% 20%, var(--noche2), var(--noche) 70%)}
  canvas{display:block;touch-action:manipulation}
  #hud{position:fixed;top:14px;left:0;right:0;text-align:center;color:var(--crema);
    font-weight:700;font-size:20px;letter-spacing:.5px;pointer-events:none;text-shadow:0 2px 8px rgba(0,0,0,.5)}
  #hint{position:fixed;bottom:16px;left:0;right:0;text-align:center;color:var(--menta);
    font-weight:500;font-size:15px;opacity:.85;pointer-events:none}
  #win{position:fixed;inset:0;display:none;flex-direction:column;align-items:center;
    justify-content:center;text-align:center;background:rgba(20,12,40,.55);backdrop-filter:blur(2px)}
  #win.show{display:flex}
  #win h1{color:var(--dorado);font-weight:800;font-size:clamp(34px,9vw,92px);line-height:1.05;
    text-shadow:0 0 30px rgba(255,200,87,.55);animation:pop .6s cubic-bezier(.2,1.6,.4,1) both}
  #win h2{color:var(--frutilla);font-weight:800;font-size:clamp(26px,7vw,64px);margin-top:6px;
    animation:pop .6s .15s cubic-bezier(.2,1.6,.4,1) both}
  #win p{color:var(--crema);font-size:clamp(15px,3vw,22px);font-weight:500;margin-top:14px;
    animation:fadeup .8s .45s both}
  #win button{margin-top:30px;border:0;cursor:pointer;font:inherit;font-weight:700;font-size:17px;
    color:var(--noche);background:var(--menta);padding:12px 28px;border-radius:999px;
    box-shadow:0 6px 0 rgba(0,0,0,.25);animation:fadeup .8s .7s both}
  #win button:active{transform:translateY(3px);box-shadow:0 3px 0 rgba(0,0,0,.25)}
  @keyframes pop{from{transform:scale(.3);opacity:0}to{transform:scale(1);opacity:1}}
  @keyframes fadeup{from{transform:translateY(14px);opacity:0}to{transform:translateY(0);opacity:1}}
  @media (prefers-reduced-motion: reduce){#win h1,#win h2,#win p,#win button{animation:none}}
</style>
</head>
<body>
<canvas id="c"></canvas>
<div id="hud">Velas: <span id="score">0</span>/5</div>
<div id="hint">Toc&aacute; la pantalla o apret&aacute; ESPACIO para disparar una vela</div>
<div id="win">
  <h1>&iexcl;FELIZ CUMPLEA&Ntilde;OS!</h1>
  <h2>&#19977;&#37070; &middot; Saburou</h2>
  <p>que se cumplan todos tus deseos</p>
  <button id="again">Jugar de nuevo</button>
</div>
<script>
(function(){
  var cv=document.getElementById('c'),cx2=cv.getContext('2d');
  var scoreEl=document.getElementById('score'),winEl=document.getElementById('win');
  var W,H,CX,CY,R;
  function resize(){
    W=cv.width=window.innerWidth; H=cv.height=window.innerHeight;
    CX=W/2; CY=H*0.34; R=Math.min(W,H)*0.21;
  }
  window.addEventListener('resize',resize); resize();

  var TAU=Math.PI*2, NSLOTS=5, TOL=0.30;
  var state,slots,rotation,speed,shot,misses,parts,floats,shake,placed,boomT;

  var sprinkles=[];
  for(var i=0;i<26;i++) sprinkles.push({a:Math.random()*TAU, r:0.25+Math.random()*0.5,
    c:['#7bdff2','#ffc857','#ffffff','#b388ff'][i%4]});

  function reset(){
    state='play'; rotation=0; speed=0.9; placed=0; boomT=0; shake=0;
    shot=null; misses=[]; parts=[]; floats=[];
    slots=[];
    for(var i=0;i<NSLOTS;i++) slots.push({a:i*TAU/NSLOTS, filled:false,
      color:['#ff5d8f','#7bdff2','#ffc857','#b388ff','#80ed99'][i]});
    scoreEl.textContent='0';
    winEl.classList.remove('show');
  }
  reset();

  function norm(a){ a=a%TAU; if(a<0)a+=TAU; return a; }
  function angDist(a,b){ var d=Math.abs(norm(a)-norm(b)); return Math.min(d,TAU-d); }

  function shoot(){
    if(state!=='play'||shot) return;
    shot={x:CX, y:H-70, v:Math.max(900,H*1.3)};
  }
  window.addEventListener('keydown',function(e){
    if(e.code==='Space'){ e.preventDefault(); shoot(); }
  });
  cv.addEventListener('pointerdown',shoot);
  document.getElementById('again').addEventListener('click',function(e){
    e.stopPropagation(); reset();
  });

  function addFloat(txt,col){ floats.push({txt:txt,col:col,y:CY+R+60,a:1}); }

  function explode(){
    state='boom'; boomT=0;
    var cols=['#ff5d8f','#ffc857','#7bdff2','#b388ff','#80ed99','#ffe8c9'];
    for(var i=0;i<220;i++){
      var an=Math.random()*TAU, sp=120+Math.random()*620;
      parts.push({x:CX,y:CY,vx:Math.cos(an)*sp,vy:Math.sin(an)*sp-150,
        s:3+Math.random()*6,c:cols[i%cols.length],rot:Math.random()*TAU,
        vr:(Math.random()-0.5)*10,a:1});
    }
  }

  function update(dt){
    if(state==='play'){
      rotation=norm(rotation+speed*dt);
      if(shot){
        shot.y-=shot.v*dt;
        var tip=shot.y-46;
        if(tip<=CY+R){
          var local=norm(Math.PI/2-rotation), best=0,bd=99;
          for(var i=0;i<NSLOTS;i++){
            var d=angDist(local,slots[i].a);
            if(d<bd){bd=d;best=i;}
          }
          if(bd<=TOL && !slots[best].filled){
            slots[best].filled=true; placed++;
            scoreEl.textContent=placed;
            speed+=0.35*(speed>0?1:-1);
            if(placed===2||placed===4) speed=-speed;
            addFloat('¡Buena!','#80ed99');
            if(placed>=NSLOTS) explode();
          }else{
            shake=0.25;
            addFloat(bd<=TOL?'¡Chocaste una vela!':'Uy, casi...','#ff5d8f');
            misses.push({x:shot.x,y:shot.y,vx:(Math.random()<0.5?-1:1)*(150+Math.random()*200),
              vy:-150,rot:0,vr:(Math.random()-0.5)*14,a:1});
          }
          shot=null;
        }
      }
    } else if(state==='boom'){
      boomT+=dt;
      if(boomT>0.9 && !winEl.classList.contains('show')) winEl.classList.add('show');
    }
    for(var i=misses.length-1;i>=0;i--){
      var m=misses[i]; m.vy+=1400*dt; m.x+=m.vx*dt; m.y+=m.vy*dt;
      m.rot+=m.vr*dt; m.a-=0.8*dt;
      if(m.a<=0||m.y>H+80) misses.splice(i,1);
    }
    for(var i=parts.length-1;i>=0;i--){
      var p=parts[i]; p.vy+=500*dt; p.x+=p.vx*dt; p.y+=p.vy*dt;
      p.rot+=p.vr*dt; p.a-=0.35*dt;
      if(p.a<=0) parts.splice(i,1);
    }
    for(var i=floats.length-1;i>=0;i--){
      floats[i].y-=40*dt; floats[i].a-=1.1*dt;
      if(floats[i].a<=0) floats.splice(i,1);
    }
    if(shake>0) shake-=dt;
  }

  function drawCandle(g,scale,lit){
    var w=10*scale,h=46*scale;
    g.fillStyle='#fff';
    g.fillRect(-w/2,-h,w,h);
    g.fillStyle='#ff5d8f';
    for(var s=0;s<3;s++) g.fillRect(-w/2,-h+6*scale+s*14*scale,w,5*scale);
    if(lit){
      var f=1+Math.sin(performance.now()/90)*0.18;
      g.fillStyle='rgba(255,159,28,0.35)';
      g.beginPath(); g.arc(0,-h-9*scale,9*scale*f,0,TAU); g.fill();
      g.fillStyle='#ffc857';
      g.beginPath(); g.ellipse(0,-h-8*scale,3.4*scale,7*scale*f,0,0,TAU); g.fill();
      g.fillStyle='#fff';
      g.beginPath(); g.ellipse(0,-h-6*scale,1.5*scale,3*scale,0,0,TAU); g.fill();
    }
  }

  function draw(){
    cx2.clearRect(0,0,W,H);
    var sx=shake>0?(Math.random()-0.5)*10:0, sy=shake>0?(Math.random()-0.5)*10:0;

    if(state!=='boom'||boomT<0.25){
      var fade=state==='boom'?Math.max(0,1-boomT*4):1;
      cx2.save();
      cx2.translate(CX+sx,CY+sy);
      cx2.globalAlpha=fade;
      cx2.scale(state==='boom'?1+boomT*2:1, state==='boom'?1+boomT*2:1);

      // resplandor y plato
      cx2.fillStyle='rgba(255,200,87,0.08)';
      cx2.beginPath(); cx2.arc(0,0,R*1.45,0,TAU); cx2.fill();
      cx2.fillStyle='#3b2a5e';
      cx2.beginPath(); cx2.arc(0,0,R*1.18,0,TAU); cx2.fill();

      cx2.rotate(rotation);

      // torta vista desde arriba
      cx2.fillStyle='#ff5d8f';
      cx2.beginPath(); cx2.arc(0,0,R,0,TAU); cx2.fill();
      cx2.fillStyle='#ff86ab';
      cx2.beginPath(); cx2.arc(0,0,R*0.86,0,TAU); cx2.fill();
      cx2.fillStyle='#ffe8c9';
      cx2.beginPath(); cx2.arc(0,0,R*0.66,0,TAU); cx2.fill();
      for(var i=0;i<sprinkles.length;i++){
        var sp=sprinkles[i];
        cx2.fillStyle=sp.c;
        cx2.save();
        cx2.rotate(sp.a); cx2.translate(R*0.66*sp.r,0); cx2.rotate(1);
        cx2.fillRect(-1.5,-4,3,8);
        cx2.restore();
      }
      cx2.fillStyle='#e63956';
      cx2.beginPath(); cx2.arc(0,0,R*0.09,0,TAU); cx2.fill();

      // portavelas y velas colocadas
      for(var i=0;i<NSLOTS;i++){
        var s=slots[i];
        cx2.save();
        cx2.rotate(s.a); cx2.translate(R*0.93,0); cx2.rotate(Math.PI/2);
        if(s.filled){
          drawCandle(cx2,1,true);
        }else{
          cx2.strokeStyle='#ffc857'; cx2.lineWidth=3;
          cx2.beginPath(); cx2.arc(0,0,9,0,TAU); cx2.stroke();
          cx2.fillStyle='rgba(255,200,87,0.25)';
          cx2.beginPath(); cx2.arc(0,0,9,0,TAU); cx2.fill();
        }
        cx2.restore();
      }
      cx2.restore();
    }

    // vela en vuelo
    if(shot){
      cx2.save(); cx2.translate(shot.x,shot.y);
      drawCandle(cx2,1.1,true);
      cx2.restore();
    }
    // lanzador (vela lista)
    if(state==='play'&&!shot){
      cx2.save(); cx2.translate(CX,H-70); cx2.globalAlpha=0.95;
      drawCandle(cx2,1.1,true);
      cx2.restore();
    }
    // velas erradas
    for(var i=0;i<misses.length;i++){
      var m=misses[i];
      cx2.save(); cx2.translate(m.x,m.y); cx2.rotate(m.rot); cx2.globalAlpha=Math.max(0,m.a);
      drawCandle(cx2,1,false);
      cx2.restore();
    }
    // confeti
    for(var i=0;i<parts.length;i++){
      var p=parts[i];
      cx2.save(); cx2.translate(p.x,p.y); cx2.rotate(p.rot); cx2.globalAlpha=Math.max(0,p.a);
      cx2.fillStyle=p.c; cx2.fillRect(-p.s/2,-p.s/2,p.s,p.s*1.6);
      cx2.restore();
    }
    // textos flotantes
    cx2.textAlign='center'; cx2.font='700 26px "Baloo 2",sans-serif';
    for(var i=0;i<floats.length;i++){
      var f=floats[i];
      cx2.globalAlpha=Math.max(0,f.a); cx2.fillStyle=f.col;
      cx2.fillText(f.txt,CX,f.y);
    }
    cx2.globalAlpha=1;
  }

  var last=performance.now();
  function loop(t){
    var dt=Math.min(0.033,(t-last)/1000); last=t;
    update(dt); draw();
    requestAnimationFrame(loop);
  }
  requestAnimationFrame(loop);
})();
</script>
</body>
</html>`

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	})
	log.Printf("Torta lista en http://localhost:%s — a soplar velas", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
