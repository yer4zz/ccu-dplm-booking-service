
create extension if not exists "uuid-ossp";
create extension if not exists "btree_gist"; 


create table public.profiles (
  id          uuid primary key references auth.users(id) on delete cascade,
  role        text not null default 'client'
              check (role in ('client', 'master', 'admin')),
  full_name   text not null,
  phone       text,
  avatar_url  text,
  created_at  timestamptz default now()
);

create or replace function public.handle_new_user()
returns trigger language plpgsql security definer as $$
begin
  insert into public.profiles (id, full_name, role)
  values (
    new.id,
    coalesce(new.raw_user_meta_data->>'full_name', 'Гость'),
    coalesce(new.raw_user_meta_data->>'role', 'client')
  );
  return new;
end;
$$;

create trigger on_auth_user_created
  after insert on auth.users
  for each row execute procedure public.handle_new_user();


create table public.masters (
  id               uuid primary key references public.profiles(id) on delete cascade,
  bio              text,
  experience_years smallint default 0,
  instagram        text,
  is_active        boolean default true
);


create table public.services (
  id           uuid primary key default uuid_generate_v4(),
  name         text not null,              
  category     text not null default 'hair'
              check (category in ('hair','beard','care','combo')),
  description  text,
  duration_min integer not null,
  price        numeric(10,2) not null,
  is_active    boolean default true,
  sort_order   smallint default 0,
  created_at   timestamptz default now()
);

create table public.master_services (
  master_id    uuid references public.masters(id) on delete cascade,
  service_id   uuid references public.services(id) on delete cascade,
  custom_price numeric(10,2),
  primary key (master_id, service_id)
);


create table public.master_schedules (
  id          uuid primary key default uuid_generate_v4(),
  master_id   uuid references public.masters(id) on delete cascade,
  day_of_week smallint not null check (day_of_week between 0 and 6),
  start_time  time not null,
  end_time    time not null,
  unique (master_id, day_of_week)
);


create table public.bookings (
  id          uuid primary key default uuid_generate_v4(),
  client_id   uuid references public.profiles(id) on delete set null,
  master_id   uuid not null references public.masters(id),
  service_id  uuid not null references public.services(id),
  starts_at   timestamptz not null,
  ends_at     timestamptz not null,
  status      text not null default 'pending'
              check (status in ('pending','confirmed','cancelled','completed','no_show')),
  price_paid  numeric(10,2), 
  notes       text,      
  created_at  timestamptz default now(),

  constraint no_overlapping_bookings exclude using gist (
    master_id with =,
    tstzrange(starts_at, ends_at, '[)') with &&
  ) where (status not in ('cancelled', 'no_show'))
);

create index idx_bookings_master_time on public.bookings (master_id, starts_at);
create index idx_bookings_client      on public.bookings (client_id, starts_at);
create index idx_bookings_status      on public.bookings (status);


create table public.reviews (
  id          uuid primary key default uuid_generate_v4(),
  booking_id  uuid unique references public.bookings(id) on delete cascade,
  rating      smallint not null check (rating between 1 and 5),
  comment     text,
  is_visible  boolean default true,
  created_at  timestamptz default now()
);


alter table public.bookings enable row level security;
alter table public.profiles enable row level security;
alter table public.reviews  enable row level security;


create policy "bookings_client_select" on public.bookings for select
  using (auth.uid() = client_id);


create policy "bookings_master_select" on public.bookings for select
  using (auth.uid() = master_id);


create policy "bookings_master_update" on public.bookings for update
  using (auth.uid() = master_id)
  with check (auth.uid() = master_id);


create policy "bookings_client_insert" on public.bookings for insert
  with check (auth.uid() = client_id);


create policy "profiles_self" on public.profiles for all
  using (auth.uid() = id);