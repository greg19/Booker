# Booker

## Opis funkcjonalności

Aplikacja będzie pozwalać na rezerwowanie terminów wizyt. Usługodawcy będą mogli
skonfigurować siatkę godzin oraz formularz, a klienci wybierać dogodne dla siebie
terminy.

## Lista funkcji

- konfiguracja dostępnych terminów
- konfiguracja dodatkowych pól do wizyty
- rezerwacja wizyt
- odwoływanie wizyt
- przesuwanie wizyt
- weryfikacja niezarejestrowanych użytkowników
- rejestracja użytkowników
- przypomnienia o wizycie
- historia wizyt
- lista przyszłych wizyt

### Panel administratora

Administrator będzie mógł skonfigurować siatkę dostępnych wizyt, jak i również
odwoływać i przesuwać dowolne wizyty. Można będzie też skonfigurować formularz,
który wypełnia klient podczas rejestracji, dodając dodatkowe pola (np. numer
telefonu, rodzaj usługi). Od niezarejestrowanych użytkowników powinna być
możliwość ustawienia wymagania przejścia dodatkowej weryfikacji (np.
potwierdzenie email albo captcha).

### Użytkownicy

Użytkownicy będą widzieć siatkę godzin z dostępnymi terminami, za pomocą której
będą mogli umówić się na wizytę. Powinni mieć również możliwość odwołania lub
przesunięcia wizyty, bez konieczności posiadania konta. Aplikacja powinna też
umożliwiać ustawienie mailowego przypomnienia o wizycie.

### Zarejestrowani użytkownicy

Zarejestrowany użytkownik będzie mógł zobaczyć listę swoich przeszłych
i przyszłych wizyt. Będzie mógł również w łatwy sposób odwołać lub przełożyć
swoje nadchodzące wizyty. Dodatkowo formularz rezerwacji może automatycznie
pobierać niektóre dane (jak imię i nazwisko albo numer telefonu).

## Implementacja

### Technologia

Wnętrze aplikacji będzie napisane w języku Go, natomiast interfejs użytkownika
będzie używał HTML, CSS i JavaScript. Jako serwis baz danych wykorzystany będzie
MySQL.

### Funkcje

W ramach projektu zostanie zaimplementowana podstawowa wersja aplikacji,
obsługująca niezarejestrowanych użytkowników. Ponadto nie planuje się
implementować weryfikacji podczas rezerwacji, ani przesuwania wizyt (ten sam
efekt można otrzymać rejestrując się na nowy termin oraz odwołując starą
wizytę).

