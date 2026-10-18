% Example Prolog program written for scc.
% Prolog is a logic programming language built on facts, rules, and backtracking.

/* Family relationships: facts */
parent(tom, bob).
parent(tom, liz).
parent(bob, ann).
parent(bob, pat).

% A grandparent relation: a rule whose body is a conjunction of two goals.
grandparent(X, Z) :- parent(X, Y), parent(Y, Z).

% Sibling relation, using \== for structural inequality.
sibling(X, Y) :- parent(Z, X), parent(Z, Y), X \== Y.

% classify/2 demonstrates arithmetic comparisons and the cut.
classify(N, negative) :- N < 0, !.
classify(N, zero) :- N =:= 0, !.
classify(N, positive) :- N > 0.

% max/3 using if-then-else syntax.
max(X, Y, Z) :- (X >= Y -> Z = X ; Z = Y).

% factorial/2, recursive arithmetic.
factorial(0, 1) :- !.
factorial(N, F) :- N > 0, N1 is N - 1, factorial(N1, F1), F is N * F1.

% adult/1 uses negation as failure.
child(bob).
child(liz).
adult(X) :- \+ child(X).

main :-
    grandparent(tom, ann),
    write('tom is a grandparent of ann'), nl,
    classify(-3, C1), write(C1), nl,
    max(4, 9, M), write(M), nl,
    factorial(5, F), write(F), nl,
    ( adult(tom) -> write('tom is an adult') ; write('tom is not an adult') ), nl.

:- initialization(main).
