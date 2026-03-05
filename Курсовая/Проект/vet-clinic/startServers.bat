@echo off
echo Starting backend...
start cmd /c "cd backend && npm run start:dev"

echo Starting frontend...
start cmd /c "cd frontend && npm run dev"

timeout /t 5 /nobreak > nul

echo Opening browser...
start http://localhost:5173

echo Done!