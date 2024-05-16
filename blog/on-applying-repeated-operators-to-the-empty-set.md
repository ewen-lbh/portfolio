---
date: 2024-05-15
title: On Applying Repeated Operators to the Empty Set
tags: [math]
mathjax: yes
---


This is a sort of "theorem" that I made up in my mind has been itching me since my years in higher math education.

## Context

Let $E$ a set and $o \in E^2 \to E$, such that $(E, o)$ forms a [monoid](https://en.wikipedia.org/wiki/Monoid). 

We then define $\mathcal O$ as the _repeated variant_ of the binary operator $o$:

$$
\begin{align*}
\mathcal O : E^\mathbb N &\to E \\
 (a_n)_n &\mapsto a_0\ o\ a_1\ o\ a_2\ o\ \ldots\ o\ a_n
\end{align*}
$$

using an infix notation for $o$, defined as you would expect: $A\ o\ B = o(A, B)$.


This is why we need $(E, o)$ to be a monoid, instead of a unital magma: we need the operator to be associative, so that the repeated application of the operator is well-defined.

Note how repeated operators have their argument in a sequence space instead of a set. This is because:

1. we need to _iterate_ over the elements, which requires a well-defined order on the elements (otherwise, we would need $o$ to be commutative and therefore $(E, o)$ to be a group instead of just a monoid)
1. We also want to be able to represent repeated operations on duplicates, which sets cannot represent.


## Theorem

For any repeated operator $\mathcal O$:

$$
\mathcal O(()) = \operatorname{unit} E
$$

where $\operatorname{unit} E$ is the unit element of the monoid $(E, o)$.

## Application to ∀ and ∃

This looks really abstract right now, but consider the two statements that we (at least in France) are told to be "self-evident":

- For any proposition $\mathcal P$, $\forall e \in \emptyset,\ \mathcal P(e)$ is true
- For any proposition $\mathcal P$, $\exists e \in \emptyset,\ \mathcal P(e)$ is false

But $\forall$ (resp. $\exists$) are just syntactic sugar for the repeated variants of the logical and $\land$ (resp. logical or $\lor$) operators. So, in fact, the two statements are equivalent to:

- For any proposition $\mathcal P$, $\bigwedge_{e \in \emptyset} \mathcal P(e)$ is true
- For any proposition $\mathcal P$, $\bigvee_{e \in \emptyset} \mathcal P(e)$ is false

Now, we can prove these statements with the preceding theorem:

1. $\bigwedge$ is the repeated variant of $\land$.
1. $(\{\top, \bot\}, \land)$ is a monoid:
    - $\land$ is a binary operation on $\{\top, \bot\}$.
    - $\land$ is associative: $a \land (b \land c) = (a \land b) \land c$.
    - $\land$ has a unit element $\top$ (as $a \land \top = a$ for any $a$) and $\top \in \{\top, \bot\}$.
1. $\operatorname{unit} \{\top, \bot\} = \top$.

Therefore, $\bigwedge(()) = \top$.

## Adapting conventional set notation

You'll have noticed how the previous statement is kind of akwardly written: we say $\mathcal O(()) = \ldots$ instead of the much more usual notations, $\mathcal O_{a \in \emptyset} a = \ldots$ or $\mathcal O_{\emptyset} = \ldots$.

This is because we decided to model the repeated operators as functions over sequences, instead of sets.

But, as long as the set is ordered, we can trivially adapt the notation:

Let $(E, \geq)$ an ordered set and $o$ a binary operator such that $(E, o)$ forms a monoid.

We define the set-to-sequence function $\operatorname{seq}$ as:

$$
\begin{align*}
\operatorname{seq} : \mathcal P(E) &\to E^\mathbb N \\
\emptyset &\mapsto () \\
\{a\} &\mapsto (a) \\
\{a\} \cup A &\mapsto \operatorname{seq}(\{ e \in A | e \lt a \}) \sqcup \{a\} \sqcup \operatorname{seq}(\{ e \in A | e \geq a \})
\end{align*}
$$

You'll notice that $\operatorname{seq}$ is basically a [quicksort](https://en.wikipedia.org/wiki/Quicksort) algorithm.

The interesting thing to note though, is that the function is _not_ bijective, as converting a sequence back to a set would require dropping duplicates. But most usages of repeated operators don't operate on duplicate elements anyways.

With that said, we can overload the notation of $\mathcal O$ to accept sets:

$$
\forall P \in \mathcal P(E), \quad \mathcal O_P := \mathcal O(\operatorname{seq}(P))
$$

Then, we can finally state:

$$
\begin{align*}
\bigwedge_\emptyset &= \top \\
\bigvee_\emptyset &= \bot
\end{align*}
$$

## Proof

The proof is actually reaaaally trivial, that's why I put "theorem" in quotes in the introduction. It's more of a way to have fun with (excessive?) formalization of simple things haha

Anyway, here is the proof.

Let $E$ a set and $o \in E^2 \to E$ such that $(E, o)$ is a monoid. Let $\mathcal O$ be the repeated variant of $o$ and $e$ the unit element of $(E, o)$.

**Proof by contradiction.**

Assume that $\mathcal O(()) \neq e$.

Then:

$$
\begin{align*}
    \mathcal O((a_1)) &= \mathcal O((a_1) \sqcup ()) \\
                      &= a_1\ o\ \mathcal O(()) \\
\end{align*}
$$

But we also have, by definition: $\mathcal O((a_1)) = a_1 = a_1\ o\ e$.

We thus have $e = \mathcal O(())$, which is a contradiction.

So, $\mathcal O(()) = e$.

## Applications

| Operation | Application |  |
| --------- | ------------------------------------------------------ |-- |
| $+$ | $\sum_\emptyset = 0$ | |
| $\cdot$ | $\prod_\emptyset = 1$ | |
| $\max$ | $\max_\emptyset = -\infty$ | |
| $\min$ | $\min_\emptyset = \infty$ | |
| $\land$ | $\forall \emptyset = \top$ | |
| $\lor$ | $\exists \emptyset = \bot$ | |
| $\cup$ | $\bigcup_\emptyset = \emptyset$ | |
| $\cap$ | $\bigcap_\emptyset = \mathbb U$ | where $\mathbb U$ is the universe set |
